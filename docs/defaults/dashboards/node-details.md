# Node Details

Prometheus node exporter sample dashboard

## Tags

* `linux`

## Panels

### Quick CPU / Mem / Disk

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| CPU Busy | Busy state of all CPU cores together | Default:<br/>Mode: absolute<br/>Level 1: 85<br/>Level 2: 95<br/><br/> |  |
| Sys Load (5m avg) | Busy state of all CPU cores together (5 min average) | Default:<br/>Mode: absolute<br/>Level 1: 85<br/>Level 2: 95<br/><br/> |  |
| Sys Load (15m avg) | Busy state of all CPU cores together (15 min average) | Default:<br/>Mode: absolute<br/>Level 1: 85<br/>Level 2: 95<br/><br/> |  |
| RAM Used | Non available RAM memory | Default:<br/>Mode: absolute<br/>Level 1: 80<br/>Level 2: 90<br/><br/> |  |
| SWAP Used | Used Swap | Default:<br/>Mode: absolute<br/>Level 1: 10<br/>Level 2: 25<br/><br/> |  |
| Root FS Used | Used Root FS | Default:<br/>Mode: absolute<br/>Level 1: 80<br/>Level 2: 90<br/><br/> |  |
| CPU Cores | Total number of CPU cores | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Uptime | System uptime | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Network IO | Shows the average network io over the last 5 minutes |  |  |
| RootFS Total | Total RootFS | Default:<br/>Mode: absolute<br/>Level 1: 70<br/>Level 2: 90<br/><br/> |  |
| RAM Total | Total RAM | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| SWAP Total | Total SWAP | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Basic CPU / Mem / Net / Disk

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| CPU Basic | Basic CPU info |  |  |
| Memory Basic | Basic memory usage |  |  |
| Network Traffic Basic | Basic network info per interface |  |  |
| Disk Space Used Basic | Disk space used of all filesystems mounted |  |  |
<!-- markdownlint-enable line-length -->

### CPU / Memory / Net / Disk

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| CPU | Shows how much time in percentage was spent on the execution of each state |  |  |
| Memory Stack | Shows how much memory use by each memory area (like Application or Swap) on selected node. Show values in stacked view. |  |  |
| Network Traffic | Show how much network traffic was processed by each interface on selected node. On graph above zero show statistic by received traffic, below zero - by transmitted traffic. |  |  |
| Disk Space Used | Show how much storage use by each mount on selected node. Show stacked values. |  |  |
| Disk IOps | Show how much disk input/output operations per second (IOPS) are used for read or write to disk on selected node. |  |  |
| I/O Usage Read / Write | Show how much disk input/output operations in bytes are used for read or write to disk on selected node. |  |  |
| I/O Usage Times | Show how much time in seconds disk spent to input/output operations on selected node. |  |  |
<!-- markdownlint-enable line-length -->

### Memory Meminfo

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Memory Active / Inactive | Show how much memory on the basis of active or inactive used on selected node.<br/>Active - Memory that has been used more recently and usually not reclaimed unless absolutely necessary.<br/>Inactive - Memory which has been less recently used.  It is more eligible to be reclaimed for other purposes |  |  |
| Memory Commited | Show how much committed memory used on selected node. |  |  |
| Memory Active / Inactive Detail | Show detailed information how much memory on the basis of active or inactive used on selected node. |  |  |
| Memory Writeback and Dirty | Show how much write back and dirty memory used on selected node. |  |  |
| Memory Shared and Mapped | Show how much shared and mapped memory used on selected node. |  |  |
| Memory Slab | Show how much memory used by Slab mechanism to reclaim and reduce fragmentation on selected node. |  |  |
| Memory Vmalloc | Show how much virtual memory (vmalloc) were allocated on selected node |  |  |
| Memory Bounce | Show how much memory was using for bounce buffers on selected node. <br/>Bounce buffers are required for devices that cannot access the full range of memory available to the CPU.<br/>An obvious example of this is when a device does not address with as many bits as the CPU, such as 32 bit devices on 64 bit architectures or recent Intel processors with PAE enabled. |  |  |
| Memory Anonymous | Show how much memory were marked as anonymous memory (or anonymous mapping) and can not be backed by a filesystem. <br/>Usually, the anonymous mappings only define virtual memory areas that the program is allowed to access |  |  |
| Memory Kernel / CPU | Show how much memory reserved and used by Kernel itself, and per CPU memory allocated dynamically by loadable modules.<br/>Kernel memory can not be reclaimed. |  |  |
| Memory HugePages Counter | Show how much page counts used by HugePages feature in Linux on selected node. memory allocated dynamically by loadable modules. |  |  |
| Memory HugePages Size | Show how much memory in bytes used by HugePages feature in Linux on selected node. memory allocated dynamically by loadable modules. |  |  |
| Memory DirectMap | Show how much DirectMap memory pages used on selected node. Show by different page size like: 4KiB, 64KiB, 1MiB, 2MiB, 4MiB, 1GiB, or 2GiB. |  |  |
| Memory Unevictable and MLocked | Show how much memory in bytes were marked as Unevictable (it means that it can not be evicted or reclaimed) or locked by system calls |  |  |
| Memory NFS | Show how much memory in NFS pages sent to the server, but not yet committed to the storage |  |  |
<!-- markdownlint-enable line-length -->

### Memory Vmstat

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Memory Pages In / Out | Show how much memory pages were used or were released by OS on selected node |  |  |
| Memory Pages Swap In / Out | Show how much memory pages were swapped or return back from swap on selected node |  |  |
| Memory Page Faults | Show how much memory page faults occurred during work on selected node. <br/>`Major page fault` means that application requested load data in memory from disk, but data still not load. Application will wait load data in memory.<br/>`Minor page fault` means that requested by application data were already loaded in memory, but were not marked as loaded. This could happen if the memory is shared by different programs and the page is already brought into memory for other programs. |  |  |
| OOM Killer | Show how much OOMKiller calls occurred on selected node. |  |  |
<!-- markdownlint-enable line-length -->

### System Timesync

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Time Syncronized Drift | Show information about time synchronization with central NTP server.   |  |  |
| Time PLL Adjust | Show the size of PLL constant which  used as the poll interval from NTP server |  |  |
| Time Syncronized Status | Show status of time synchronization selected node with NTP server. |  |  |
| Time Misc | Show various metrics about time synchronization which can not be group.: ticks between seconds, International Atomic Time (TAI) offset |  |  |
<!-- markdownlint-enable line-length -->

### System Processes

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Processes Status | Show how much system processes work in running or block (by I/O operations) states |  |  |
| Processes State | Show count of processes by states, where states:<br/>* `D` uninterruptible sleep (usually IO) <br/>* `I` Idle kernel thread <br/>* `R` running or runnable (on run queue) <br/>* `S` interruptible sleep (waiting for an event to complete) <br/>* `T` stopped by job control signal <br/>* `t` stopped by debugger during the tracing <br/>* `W` paging (not valid since the 2.6.xx kernel) <br/>* `X` dead (should never be seen) <br/>* `Z` defunct ("zombie") process, terminated but not reaped by its parent |  |  |
| Processes  Forks | Show how much forks were created by running processes on selected node |  |  |
| Processes Memory | Show how much memory used by processes on selected node |  |  |
| Process schedule stats Running / Waiting | Show how much processes spend time in running state or in waiting processor time. Can be used to spot not only if you have more processes to be run than CPU time available to handle them, but also indicate how many such processes there on average. |  |  |
| PIDs Number and Limit | Show how much pids used and limit by processes / pids on selected node |  |  |
| Threads Number and Limit | Show how much threads used and limit by threads count on selected node |  |  |
<!-- markdownlint-enable line-length -->

### System Misc

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Context Switches / Interrupts | Show how much context switches or interrupts operations occurred in system on selected node |  |  |
| System Load | Show load average of system for last `1m`, `5m`, `15m`.  Load average calculated by Linux kernel based on CPU load, I/O operation and count/state of processes.  |  |  |
| Interrupts Detail | Show detailed information about system interrupts if metrics collection interrupts enabled. Disabled by default and should be enable in `node_exporter`. |  |  |
| Schedule timeslices executed by each cpu | Show the number of slices of time that were used to schedule process to running  |  |  |
| Entropy | Show the size of available entropy in bytes. Entropy use in such random sources as `/dev/random` or `/dev/urandom`. In its turn these sources use in various random generators in various program languages. Depends on entropy size you can have not enough data to correct work random. |  |  |
| CPU time spent in user and system contexts | Show how much time in seconds CPU spent in user and system contexts |  |  |
| File Descriptors | Show mow much file descriptors use and limit by them on selected node. File descriptors use to create processes (even create SSH sessions). So not enough count of file descriptors can affect work of system.   |  |  |
<!-- markdownlint-enable line-length -->

### Hardware Misc

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Hardware temperature monitor | Show various temperatures of hardware on selected node. These metrics are no always available. Usually virtual machines does not expose these metrics. |  |  |
| Throttle cooling device | Show throttle cooling device. These metrics are no always available. Usually virtual machines does not expose these metrics. |  |  |
| Power supply | Show state of power supply. These metrics are no always available. Usually virtual machines does not expose these metrics. |  |  |
<!-- markdownlint-enable line-length -->

### Systemd

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Systemd Sockets | Show how much connections exists to systemd sockets if metrics collection enabled. Disabled by default and to enable need add special flag to `node_exporter`. |  |  |
| Systemd Units State | Show systemd units count by state (`activating`, `active`, `failed` and so on) on selected node. Disabled by default and to enable need add special flag to `node_exporter`. |  |  |
<!-- markdownlint-enable line-length -->

### Storage Disk

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Disk IOps Completed | Show how much reads/writes in second were completed by device on selected node |  |  |
| Disk R/W Data | Show how much data were read/write in bytes per second by device on selected node |  |  |
| Disk R/W Time | Show how much time were spent to read/write on disk by devices on selected node |  |  |
| Disk IOs Weighted | Show weighted time in seconds spent doing I/O operations. This value is incremented at each I/O start, I/O completion, I/O merge, or read of these stats by the number of I/Os in progress times the number of seconds spent doing I/O since the<br/>last update of this value. This can provide an easy measure of both I/O completion time and the backlog that may be accumulating. |  |  |
| Disk R/W Merged | Show the total number of reads and writes merged |  |  |
| Time Spent Doing I/Os | Show how much time in seconds were spent on doing I/O operations |  |  |
| Disk IOs Current in Progress | Show how much I/O operations per second were in progress on moment collect metric value on selected node |  |  |
| Disk IOps Discards completed / merged | Show the total number of discards completed successfully and discards merged. |  |  |
<!-- markdownlint-enable line-length -->

### Storage Filesystem

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Filesystem space available | Show how much space used and available by volume mounts on selected node |  |  |
| File Nodes Free | Show how much inodes available to use on selected node |  |  |
| File Descriptor | Show how much file descriptors were used and limit of file descriptors count on selected node |  |  |
| File Nodes Size | Show how much inodes permitted to use on selected node |  |  |
| Filesystem in ReadOnly / Error | Show filesystems which in read only state and how much problem occurs during getting information for the filesystem via `statfs`. |  |  |
<!-- markdownlint-enable line-length -->

### Network Traffic

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Network Traffic by Packets | Show how much packages were received/transmitted via interface in second on selected node  |  |  |
| Network Traffic Errors | Show how much errors occurred during received/transferred packages on selected node  |  |  |
| Network Traffic Drop | Show how much packages were dropped during received / transferred by network interfaces on selected node |  |  |
| Network Traffic Compressed | Show how much packages were compressed during received / transferred by network interfaces on selected node |  |  |
| Network Traffic Multicast | Show how much packages were send with using multicast by network interfaces on selected node |  |  |
| Network Traffic Fifo | Show network device statistic receive and transmitted fifo. |  |  |
| Network Traffic Frame | Show network device statistic receive frame |  |  |
| Network Traffic Carrier | Show network device statistic transmit carrier. |  |  |
| Network Traffic Colls | Show network device statistic transmit colls |  |  |
| NF Contrack | Show number of currently allocated flow entries for connection tracking and maximum size of connection tracking table |  |  |
| ARP Entries | Show the number of ARP Entries |  |  |
| MTU | Show the maximum transmission unit (MTU) size of one package by network interface on selected node |  |  |
| Speed | Show network interface speed in bytes per second on selected node |  |  |
| Queue Length | Show how much packages were in queue to transmit per network interface on selected node |  |  |
| Softnet Packets | Show number of processed packages |  |  |
| Softnet Out of Quota | Show number of times processing packets ran out of quota |  |  |
| Network Operational Status | Show network interface operational and physical link statuses |  |  |
<!-- markdownlint-enable line-length -->

### Network Sockstat

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Sockstat TCP | Show how much of TCP sockets in state `alloc`, `inuse`, `mem`and so on for all states |  |  |
| Sockstat UDP | Show how much of UDP sockets in state `inuse`, `mem` |  |  |
| Sockstat Used | Show number of IPv4 sockets in use |  |  |
| Sockstat Memory Size | Show number of TCP and UDP sockets in state mem_bytes |  |  |
| Sockstat FRAG / RAW | Show number of FRAG and RAW sockets in states `inuse` and `memory` |  |  |
<!-- markdownlint-enable line-length -->

### Network Netstat

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Netstat IP In / Out Octets | Show network statistic by input and output octets  |  |  |
| Netstat IP Forwarding | Show network statistic for IP forwarding |  |  |
| ICMP In / Out | Show number of ICMP messages which were received or were tried to send   |  |  |
| ICMP Errors | Show number of messages which the entity received but determined as having ICMP-specific errors (bad ICMP checksums, bad length, etc.) |  |  |
| UDP In / Out | Show how much UDP datagrams were received and transmitted  |  |  |
| UDP Errors | Show how much UDP errors occurred during send datagrams by error code |  |  |
| TCP In / Out | Sow how much TCP segments were received and transmitted  |  |  |
| TCP Errors | Show how much TCP errors occurred during send segments by error code |  |  |
| TCP Connections | Show number of TCP connections for which the current state is either ESTABLISHED or CLOSE- WAIT |  |  |
| TCP SynCookie | Show number of TCP SYN cookies received or transmitted  |  |  |
| TCP Direct Transition | Show number of TCP connections that have made a direct transition to the SYN-SENT state from the CLOSED state or to the SYN-RCVD state from the LISTEN state |  |  |
<!-- markdownlint-enable line-length -->

### Node Exporter

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Node Exporter Scrape Time | Show duration in seconds of a collector scrape in node_exporter |  |  |
| Node Exporter Scrape | Show state of scrape of collector, where:<br/>* `1` - succeeded<br/>* `0` - failed or disabled |  |  |
<!-- markdownlint-enable line-length -->
