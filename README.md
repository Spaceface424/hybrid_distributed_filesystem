# CS425_MP3_G14

UIUC CS 425 MP3 Repository 

Group Number 14

Ryan Wu (rwu29), Oliver Rogalski (oliverr3)

To run the MP, clone this repository into all of the virtual machines that will be used. Run main.go on the introducer process without the ip then provide the introducer IP to any joining processes. Use the -v flag to enable verbose output to stdout.

Example:
```
go run swim/main.go [-i introducer_ip] -v
```

To run commands, just input the command into the terminal.

Example:
```
create local_file_name hydfs_file_name
```

MP2 Commands:
- list_mem - Prints the membership list of the node, including the hash ring ID
- list_self - Prints the member's own ID
- leave - Makes the member exit and leave the network
- enable/disable_sus - Enables or disables the suspicion protocol on all machines
- status_sus - Prints whether the suspicion protocl is enabled on a member
- drop_rate - Sets all machines' drop rate to the input_argument
- list_sus - Prints all suspicious members in each machine's membership list

MP3 Commands: 
- create [localfilename] [hydfsfilename] - Creates an HYDFS file on the network
- get [hydfsfilename] [localfilename] - Reads an HYDFS file to localfilename
- append [localfilename] [hydfsfilename] - Appends localfilename to an HYDFS file
- merge [hydfsfilename] - Synchronizes all replicas of an HYDFS file
- ls [hydfsfilename] - Lists all nodes with the HYDFS file
- store - Lists all files at the current node
- getfromreplica [vmaddress] [hydfsfilename] [localfilename] - Reads an HYDFS file from a specified node to localfilename
- list_mem - Prints the membership list of the node, including the hash ring ID
- multiappend [hydfsfilename] [vmi ... vmj] [localfilenamei ... localfilenamej] - Launches appends from multiple nodes 