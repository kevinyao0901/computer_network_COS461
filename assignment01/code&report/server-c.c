/*****************************************************************************
 * server-c.c                                                                 
 * Name:
 * NetId:
 *****************************************************************************/

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <errno.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netdb.h>
#include <netinet/in.h>
#include <errno.h>

#define QUEUE_LENGTH 10
#define RECV_BUFFER_SIZE 2048

/* TODO: server()
 * Open socket and wait for client to connect
 * Print received message to stdout
 * Return 0 on success, non-zero on failure
*/

/*void sigchld_handler(int s)
{
    // waitpid() might overwrite errno, so we save and restore it:
    int saved_errno = errno;

    while (waitpid(-1, NULL, WNOHANG) > 0);

    errno = saved_errno;
}*/

int server(char *server_port) {
    int sockfd, new_fd;
    struct addrinfo hints, * servinfo, * p;
    struct sockaddr_storage their_addr;
    socklen_t sin_size;
    int yes = 1;
    char s[INET6_ADDRSTRLEN];
    int rv;
    int text_length;

    memset(&hints, 0, sizeof hints);
    hints.ai_family = AF_UNSPEC;
    hints.ai_socktype = SOCK_STREAM;
    hints.ai_flags = AI_PASSIVE;

    if ((rv = getaddrinfo(NULL, server_port, &hints, &servinfo)) != 0)
    {
        fprintf(stderr, "getaddrinfo:%s\n", gai_strerror(rv));
        return 1;
    }

    for (p = servinfo; p != NULL; p = p->ai_next)
    {
        if ((sockfd = socket(p->ai_family, p->ai_socktype,p->ai_protocol))==-1)
        {
            perror("server:socket");
            continue;
        }

        if (setsockopt(sockfd, SOL_SOCKET, SO_REUSEADDR, &yes, sizeof(int)) == -1)
        {
            perror("setsocket");
            exit(1);
        }

        if (bind(sockfd, p->ai_addr, p->ai_addrlen) == -1)
        {
            close(sockfd);
            perror("sever:bind");
            continue;
        }

        break;
    }

    freeaddrinfo(servinfo);

    if (p == NULL)
    {
        fprintf(stderr, "server:failed to bind\n");
        exit(1);
    }

    if (listen(sockfd, QUEUE_LENGTH) == -1)
    {
        perror("listen:");
        exit(1);
    }

    /*sa.sa_handler = sigchld_handler;
    sigemptyset(&sa.sa_mark);
    sa.sa_flag = SA_RESTART;
    if (sigaction(SIGCHLD, &sa, NULL) == -1)
    {
        perror("sigaction");
        exit(1);
    }

    printf("sever:waiting for connnections...\n");*/

    while (1)
    {
        sin_size = sizeof their_addr;
        new_fd = accept(sockfd, (struct sockaddr*)&their_addr, &sin_size);
        if (new_fd == -1)
        {
            perror("accept");
            continue;
        }

        char buf[RECV_BUFFER_SIZE];
        text_length = recv(new_fd, buf, RECV_BUFFER_SIZE, 0);
        while (text_length > 0)
        {
            fwrite(buf, text_length, 1, stdout);
            fflush(stdout);
            text_length = recv(new_fd, buf, RECV_BUFFER_SIZE, 0);

        }

        close(new_fd);
    }
    close(sockfd);
    return 0;
}

/*
 * main():
 * Parse command-line arguments and call server function
*/
int main(int argc, char **argv) {
  char *server_port;

  if (argc != 2) {
    fprintf(stderr, "Usage: ./server-c [server port]\n");
    exit(EXIT_FAILURE);
  }

  server_port = argv[1];
  return server(server_port);
}
