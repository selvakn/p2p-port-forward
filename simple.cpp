// Comprehensive stress test for socket-like API

#include <stdio.h>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <string.h>

#include <netdb.h>

#include "libzt.h"

#define NETWORK_ID "8056c2e21c000001"
#define PORT 7878
#define BUF_SIZE 2000


void attach(int from_fd, int to_fd) {
    char buffer[BUF_SIZE];

    while (true) {
        size_t readLength = read(from_fd, buffer, BUF_SIZE);
        write(to_fd, buffer, readLength);
    }
}


void *receive_messages(void *sockPtr) {
    attach(*(int *) sockPtr, 1);

    return NULL;
}


int listen_and_accept(int sockfd) {
    struct sockaddr_in6 serv_addr, cli_addr;

    serv_addr.sin6_flowinfo = 0;
    serv_addr.sin6_family = AF_INET6;
    serv_addr.sin6_addr = in6addr_any;
    serv_addr.sin6_port = htons(PORT);

    if (zts_bind(sockfd, (struct sockaddr *) &serv_addr, sizeof(serv_addr)) < 0)
        printf("ERROR on binding");
    printf("Bind Complete\n");

    zts_listen(sockfd, 1);
    printf("Listening\n");

    int cli_addr_len = sizeof(cli_addr);

    // accept
    int newsockfd = zts_accept(sockfd, (struct sockaddr *) &cli_addr, (socklen_t *) &cli_addr_len);
    if (newsockfd < 0)
        printf("ERROR on accept");

    char client_addr_ipv6[100];
    inet_ntop(AF_INET6, &(cli_addr.sin6_addr), client_addr_ipv6, 100);
    printf("Incoming connection from client having IPv6 address: %s\n", client_addr_ipv6);
    return newsockfd;
}

void receive_messages_in_background(int &newsockfd) {
    pthread_t rThread;
    int ret = pthread_create(&rThread, NULL, receive_messages, &newsockfd);
    if (ret != 0) {
        printf("ERROR: Return Code from pthread_create() is %d\n", ret);
    }
}

int main(int argc, char *argv[]) {
    zts_simple_start("./zt", NETWORK_ID);

    char id[ZT_ID_LEN + 1];
    zts_get_device_id(id);
    printf("id = %s\n", id);

    char homePath[ZT_HOME_PATH_MAX_LEN + 1];
    zts_get_homepath(homePath, ZT_HOME_PATH_MAX_LEN);
    printf("homePath = %s\n", homePath);

    char ipv4[ZT_MAX_IPADDR_LEN];
    char ipv6[ZT_MAX_IPADDR_LEN];
    zts_get_ipv4_address((char *) NETWORK_ID, ipv4, ZT_MAX_IPADDR_LEN);
    printf("ipv4 = %s\n", ipv4);

    zts_get_ipv6_address((char *) NETWORK_ID, ipv6, ZT_MAX_IPADDR_LEN);
    printf("ipv6 = %s\n", ipv6);

    printf("peer_count = %lu\n", zts_get_peer_count());

    int sockfd;
    if ((sockfd = zts_socket(AF_INET6, SOCK_STREAM, 0)) < 0) {
        fprintf(stderr, "error in opening socket\n");
    }
    printf("sockfd = %d\n", sockfd);

    if (argv[1]) {
        struct hostent *server = gethostbyname2(argv[1], AF_INET6);
        struct sockaddr_in6 serv_addr;
        memset((char *) &serv_addr, 0, sizeof(serv_addr));
        serv_addr.sin6_flowinfo = 0;
        serv_addr.sin6_family = AF_INET6;
        memmove((char *) &serv_addr.sin6_addr.s6_addr, (char *) server->h_addr, server->h_length);
        serv_addr.sin6_port = htons(PORT);
        int err;

        if ((err = zts_connect(sockfd, (struct sockaddr *) &serv_addr, sizeof(serv_addr))) < 0) {
            printf("error connecting to remote host (%d)\n", err);
            return -1;
        }

        receive_messages_in_background(sockfd);

        attach(0, sockfd);

    } else {
        int newsockfd = listen_and_accept(sockfd);

        receive_messages_in_background(newsockfd);

        attach(0, newsockfd);


        sleep(2);
        zts_close(newsockfd);
    }


    sleep(2);
    zts_close(sockfd);

    zts_stop();
    return 0;
}