flags= -Wall -O2 -std=c2x 
ldflags=  
tree: tree.o
	cc ${flags} $^ -o $@ ${ldflags}
tree.o: tree.c
	cc ${flags} -c $^
clean:
	rm -rf *.o tree
all: tree
