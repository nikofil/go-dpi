#include <pcap.h>
#include <libndpi-1.8.0/libndpi/ndpi_main.h>

extern int ndpi_initialize();
extern void ndpi_destroy(void);
extern int pcap_packet_callback(const struct pcap_pkthdr*, const u_char*);
