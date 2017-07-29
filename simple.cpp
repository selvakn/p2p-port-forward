// Comprehensive stress test for socket-like API

#include <stdio.h>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <string.h>

#include <netdb.h>

#include "libzt.h"

#define NETWORK_ID "8056c2e21c000001"

int main(int argc, char *argv[]) {
    zts_start("./zt");

    char id[ZT_ID_LEN + 1];
    zts_get_device_id(id);
    printf("id = %s\n", id);

    char homePath[ZT_HOME_PATH_MAX_LEN + 1];
    zts_get_homepath(homePath, ZT_HOME_PATH_MAX_LEN);
    printf("homePath = %s\n", homePath);

    // Wait for ZeroTier service to start
    while (!zts_running()) {
        printf("wating for service to start\n");
        sleep(1);
    }

    zts_join((char *) NETWORK_ID);

    while (!zts_has_ipv6_address((char *) NETWORK_ID)) {
        printf("waiting for service to issue an address\n");
        sleep(1);
    }

    while (!zts_has_ipv4_address((char *) NETWORK_ID)) {
        printf("waiting for service to issue an address\n");
        sleep(1);
    }

    char ipv4[ZT_MAX_IPADDR_LEN];
    char ipv6[ZT_MAX_IPADDR_LEN];
    zts_get_ipv4_address((char *) NETWORK_ID, ipv4, ZT_MAX_IPADDR_LEN);
    printf("ipv4 = %s\n", ipv4);

    zts_get_ipv6_address((char *) NETWORK_ID, ipv6, ZT_MAX_IPADDR_LEN);
    printf("ipv6 = %s\n", ipv6);

    printf("peer_count = %lu\n", zts_get_peer_count());

    int err;
    int sockfd;
    int port = 7878;

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
        serv_addr.sin6_port = htons(port);
        if ((err = zts_connect(sockfd, (struct sockaddr *) &serv_addr, sizeof(serv_addr))) < 0) {
            printf("error connecting to remote host (%d)\n", err);
            return -1;
        }

        zts_write(sockfd, "hello world", 11);
        sleep(3);
        zts_close(sockfd);

    } else {
        struct sockaddr_in6 serv_addr, cli_addr;

        serv_addr.sin6_flowinfo = 0;
        serv_addr.sin6_family = AF_INET6;
        serv_addr.sin6_addr = in6addr_any;
        serv_addr.sin6_port = htons(port);

        if (zts_bind(sockfd, (struct sockaddr *) &serv_addr, sizeof(serv_addr)) < 0)
            printf("ERROR on binding");

        zts_listen(sockfd, 1);

        int clilen = sizeof(cli_addr);

        // accept
        int newsockfd = zts_accept(sockfd, (struct sockaddr *) &cli_addr, (socklen_t *) &clilen);
        if (newsockfd < 0)
            printf("ERROR on accept");

        char client_addr_ipv6[100];
        inet_ntop(AF_INET6, &(cli_addr.sin6_addr), client_addr_ipv6, 100);
        printf("Incoming connection from client having IPv6 address: %s\n", client_addr_ipv6);


        printf("reading from buffer\n");
        char newbuf[32];
        memset(newbuf, 0, 32);
        read(newsockfd, newbuf, 20);
        printf("newbuf = %s\n", newbuf);

        sleep(2);
        zts_close(newsockfd);
        sleep(2);
        zts_close(sockfd);
    }


    zts_stop();
    return 0;
}
