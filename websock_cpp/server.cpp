#include <iostream>
#include <vector>
#include <cstring>
#include <cassert>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/ip.h>
#include <errno.h>

const size_t k_max_msg = 4096;

static void die(const char *msg) {
    perror(msg);
    exit(EXIT_FAILURE);
}

static ssize_t read_full(int fd, std::vector<char>& buf, size_t n) {
    size_t total_read = 0;
    while (total_read < n) {
        ssize_t rv = read(fd, buf.data() + total_read, n - total_read);
        if (rv < 0) {
            if (errno == EINTR) // Interrupted, try again
                continue;
            return -1; // Error or unexpected EOF
        }
        if (rv == 0)
            return total_read; // EOF
        total_read += rv;
    }
    return total_read;
}

static ssize_t write_full(int fd, const std::vector<char>& buf, size_t n) {
    size_t total_written = 0;
    while (total_written < n) {
        ssize_t rv = write(fd, buf.data() + total_written, n - total_written);
        if (rv < 0) {
            if (errno == EINTR) // Interrupted, try again
                continue;
            return -1; // Error
        }
        total_written += rv;
    }
    return total_written;
}

static int32_t one_request(int connfd) {
    // 4 bytes header
    std::vector<char> rbuf(4 + k_max_msg + 1);
    int32_t err = read_full(connfd, rbuf, 4);
    if (err <= 0) {
        if (err == 0)
            std::cerr << "EOF" << std::endl;
        else
            perror("read() error");
        return err;
    }

    uint32_t len = 0;
    memcpy(&len, rbuf.data(), 4);  // Assume little endian
    if (len > k_max_msg) {
        std::cerr << "Too long" << std::endl;
        return -1;
    }

    // Request body
    err = read_full(connfd, rbuf, len);
    if (err <= 0) {
        perror("read() error");
        return err;
    }

    // Do something
    rbuf[4 + len] = '\0';
    std::cout << "Client says: " << (rbuf.data() + 4) << std::endl;

    // Reply using the same protocol
    const char reply[] = "world";
    std::vector<char> wbuf(4 + sizeof(reply));
    len = strlen(reply);
    memcpy(wbuf.data(), &len, 4);
    memcpy(wbuf.data() + 4, reply, len);
    return write_full(connfd, wbuf, 4 + len);
}

int main() {
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd < 0) {
        die("socket()");
    }

    // Set socket options
    int val = 1;
    if (setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, &val, sizeof(val)) < 0) {
        die("setsockopt()");
    }

    // Bind
    struct sockaddr_in addr = {};
    addr.sin_family = AF_INET;
    addr.sin_port = htons(1234);
    addr.sin_addr.s_addr = htonl(INADDR_ANY); // Wildcard address 0.0.0.0
    if (bind(fd, reinterpret_cast<const sockaddr*>(&addr), sizeof(addr)) < 0) {
        die("bind()");
    }

    // Listen
    if (listen(fd, SOMAXCONN) < 0) {
        die("listen()");
    }

    while (true) {
        // Accept
        struct sockaddr_in client_addr = {};
        socklen_t socklen = sizeof(client_addr);
        int connfd = accept(fd, reinterpret_cast<struct sockaddr*>(&client_addr), &socklen);
        if (connfd < 0) {
            perror("accept()");
            continue;   // Error
        }

        // Serve one-request at a time
        while (true) {
            int err = one_request(connfd);
            if (err <= 0) {
                break;
            }
        }
        close(connfd);
    }

    return 0;
}
