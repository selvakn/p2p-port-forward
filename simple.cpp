// Comprehensive stress test for socket-like API

#include <stdio.h>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <arpa/inet.h>
#include <string.h>

#include <netinet/in.h>
#include <netdb.h>

#include "libzt.h"

int main() {
    char *nwid = (char *) "8056c2e21c000001";

    // Spawns a couple threads to support ZeroTier core, userspace network stack, and generates ID in ./zt
    zts_start("./zt");

    // Print the device/app ID (this is also the ID you'd see in ZeroTier Central)
    char id[ZT_ID_LEN + 1];
    zts_get_device_id(id);
    printf("id = %s\n", id);

    // Get the home path of this ZeroTier instance, where we store identity keys, conf files, etc
    char homePath[ZT_HOME_PATH_MAX_LEN + 1];
    zts_get_homepath(homePath, ZT_HOME_PATH_MAX_LEN);
    printf("homePath = %s\n", homePath);

    // Wait for ZeroTier service to start
    while (!zts_running()) {
        printf("wating for service to start\n");
        sleep(1);
    }

    // Join a network
    zts_join(nwid);

    // Wait for ZeroTier service to issue an address to the device on the given network
    while (!zts_has_ipv6_address(nwid)) {
        printf("waiting for service to issue an address\n");
        sleep(1);
    }

    while (!zts_has_ipv4_address(nwid)) {
        printf("waiting for service to issue an address\n");
        sleep(1);
    }

    // Get the ipv4 address assigned for this network
    char ipv4[ZT_MAX_IPADDR_LEN];
    char ipv6[ZT_MAX_IPADDR_LEN];
    zts_get_ipv4_address(nwid, ipv4, ZT_MAX_IPADDR_LEN);
    printf("ipv4 = %s\n", ipv4);

    zts_get_ipv6_address(nwid, ipv6, ZT_MAX_IPADDR_LEN);
    printf("ipv6 = %s\n", ipv6);

    printf("peer_count = %lu\n", zts_get_peer_count());

    // Begin Socket API calls

    int err;
    int sockfd;
    int port = 7878;
    struct sockaddr_in addr;

    // // socket()
    // if ((sockfd = zts_socket(AF_INET, SOCK_STREAM, 0)) < 0)
    //     printf("error creating ZeroTier socket");
    // else
    //     printf("sockfd = %d\n", sockfd);


    // connect() IPv6
    if (true) {
        if ((sockfd = zts_socket(AF_INET6, SOCK_STREAM, 0)) < 0) {
            fprintf(stderr, "error in opening socket\n");
        }
        printf("sockfd = %d\n", sockfd);

        struct hostent *server = gethostbyname2("fd80:56c2:e21c::199:93c5:c800:80a6", AF_INET6);
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

    }
    // connect() IPv4
    if (false) {
        addr.sin_addr.s_addr = inet_addr("28.87.3.123");
        addr.sin_family = AF_INET;
        addr.sin_port = htons(port);
        if ((err = zts_connect(sockfd, (const struct sockaddr *) &addr, sizeof(addr))) < 0) {
            printf("error connecting to remote host (%d)\n", err);
            return -1;
        }

        zts_write(sockfd, "hello", 5);
        sleep(3);
        zts_close(sockfd);
    }
    // bind() ipv4
    if (false) {
        // addr.sin_addr.s_addr = INADDR_ANY; // TODO: Requires significant socket multiplexer work
        addr.sin_addr.s_addr = inet_addr("10.9.9.40");
        // addr.sin_addr.s_addr = htons(INADDR_ANY);
        addr.sin_family = AF_INET;
        addr.sin_port = htons(port);
        if ((err = zts_bind(sockfd, (const struct sockaddr *) &addr, sizeof(addr))) < 0) {
            printf("error binding to interface (%d)\n", err);
            return -1;
        }
        zts_listen(sockfd, 1);
        struct sockaddr_in client;
        int c = sizeof(struct sockaddr_in);

        int accept_fd = zts_accept(sockfd, (struct sockaddr *) &client, (socklen_t *) &c);

        printf("reading from buffer\n");
        char newbuf[32];
        memset(newbuf, 0, 32);
        read(accept_fd, newbuf, 20);
        printf("newbuf = %s\n", newbuf);
    }

    if (false) {
        struct sockaddr_in6 me;
        if ((sockfd = zts_socket(AF_INET6, SOCK_DGRAM, IPPROTO_UDP)) < 0) {
            fprintf(stderr, "error in opening socket\n");
        }

        memset(&me, 0, sizeof(me));
        me.sin6_family = AF_INET6;
        me.sin6_port = htons(port);
//        me.sin6_addr = in6addr_any;
////        me.sin6_flowinfo = 0;
//
//
////        struct hostent *server;
////        server = gethostbyname2("::", AF_INET6);
////        memmove((char *) &me.sin6_addr.s6_addr, (char *) server->h_addr, server->h_length);
//
//        if ((err = zts_bind(sockfd, (struct sockaddr *) &me, sizeof(me))) < 0) {
//            printf("error binding to interface (%d)\n", err);
//        }
//
//        printf("Port: %d\n", ntohs(me.sin6_port));
//
//
//        int size = 1;
//        char newbuf[size];
//
//        struct sockaddr_in6 client_addr;
//        socklen_t client_length = sizeof(client_addr);
//
//        printf("reading from buffer\n");
//
//        memset(newbuf, 0, size);
//        if(recvfrom(sockfd, (void *)newbuf, size, 0, (struct sockaddr *)&client_addr, &client_length) < 0){
//          printf("error reading data\n");
//        }
//        // read(accept_fd, newbuf, 20);
//        printf("newbuf = %s\n", newbuf);
    }

    if(true) {

    }


    // End Socket API calls


    while (1) {
        sleep(1);
    }

    // Stop service, delete tap interfaces, and network stack
    zts_stop();
    return 0;
}
