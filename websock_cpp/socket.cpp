#include <bits/stdc++.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <errno.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/ip.h>

int main()
{
    int sockfd = socket(AF_INET,SOCK_STREAM, 0);
    int val = 1;
    setsockopt(sockfd,SOL_SOCKET,SO_REUSEADDR,&val,sizeof(val));
    struct sockaddr_in addr = {};
    addr.sin_family = AF_INET;
    addr.sin_port = ntohs(1234);
    addr.sin_addr.s_addr = ntohl(0); //wildcard address
    //binding the socket to the wildcard address, to listen to all interfaces
    int rv = bind(sockfd, (const sockaddr *)(&addr), sizeof(addr));

    
    
}