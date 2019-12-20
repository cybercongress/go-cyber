# Decentralization must be decentralized

January 3, 2019, we've launched first public testnet Euler-3. Since this time we have 3 relaunches and much more we'll has in the future. Thanks to our testers and validators we're finding and fixing new bugs every day. But now one fundamental and critical bug is not fixed yet. Currently, we have just 2 seed nodes and they able to upload data and provide connection to other nodes. Unfortunately, this is not about decentralization.

An obvious problem of decentralization is that no entity has a global vision of the system, and there is no central authority to direct nodes in making optimal decisions with regard to software updates, routing, or solving consensus. This makes the availability of a decentralized network more difficult to maintain, a factor significant enough to contribute to the failure of a system.

By the way, a huge part of disconnections and, as result, validators jailing happens by this reason.

Cyberd cli can’t automatically configure your router to open port `26656`, you will need to manually configure your router. We’ve can't make the following instructions to cover all router models; if you need specific help with your router, please ask for help on our [devChat](https://t.me/fuckgoogle).

Enabling inbound connections requires two steps:

1. Giving your computer a static (unchanging) internal IP address by configuring the Dynamic Host Configuration Protocol (DHCP) on your router.

2. Forwarding inbound connections from the Internet through your router to your computer where cyberd container can process them.

3. Editing cyberd configuration file.  

## Configuring DHCP

In order for your router to direct incoming port `26656` connections to your computer, it needs to know your computer’s internal IP address. However, routers usually give computers dynamic IP addresses that change frequently, so we need to ensure your router always gives your computer the same internal IP address.

Start by logging into your router’s administration interface. Most routers can be configured using one of the following URLs, so keep clicking links until you find one that works. If none work, consult your router’s manual.

```py
    http://192.168.0.1 (some Linksys/Cisco models)
    http://192.168.1.1 (some D-Link/Netgear models)
    http://192.168.2.1 (some Belkin/SMC models)
    http://192.168.123.254 (some US Robotics models)
    http://10.0.1.1 (some Apple models)
```

Upon connecting, you will probably be prompted for a username and password. If you configured a password, enter it now. If not, the Router Passwords site provides a database of known default username and password pairs.

After logging in, you want to search your router’s menus for options related to DHCP, the Dynamic Host Configuration Protocol. These options may also be called Address Reservation.

In the reservation configuration, some routers will display a list of computers and devices currently connected to your network, and then let you select a device to make its current IP address permanent.

If that’s the case, find the computer running cyberd container in the list, select it, and add it to the list of reserved addresses. Make a note of its current IP address—we’ll use the address in the next section.

Other routers require a more manual configuration. For these routers, you will need to look up the fixed address (MAC address) for your computer’s network card and add it to the list.

Open a terminal and type ifconfig. Find the result that best matches your connection—a result starting with wlan indicates a wireless connection. Find the field that starts with HWaddr and copy the immediately following field that looks like `01:23:45:67:89:ab`. Use that value in the instructions below.

Once you have the MAC address, you can fill it into to your router’s manual DHCP assignment table. Also, choose an IP address and make a note of it for the instructions in the next subsection. After entering this information, click the Add or Save button.

Then reboot your computer to ensure it gets assigned the address you selected and proceed to the Port Forwarding instructions below.

## Port Forwarding

For this step, you need to know the local IP address of the computer running cyberd container. You should have this information from configuring the DHCP assignment table in the subsection above.

Login to your router using the same steps described near the top of the DHCP subsection. Look for an option called Port Forwarding, Port Assignment, or anything with “Port” in its name. On some routers, this option is buried in an Applications & Gaming menu.

The port forwarding settings should allow you to map an external port on your router to the “internal port” of a device on your network.

Both the external port and the internal port should be `26656` for cyberd container.

Make sure the IP address you enter is the same one you configured in the previous subsection.

After filling in the details for the mapping, save the entry. You should not need to restart anything. Just ask us in [devChat](https://t.me/fuckgoogle) about successful connection.  

If you still can’t connect and you use a firewall, you probably need to change your firewall settings. Ubuntu comes with its firewall disabled by default, but if you have enabled it, see the Ubuntu [wiki page](https://help.ubuntu.com/community/Gufw) for information about adding port forwarding rules.

If something else went wrong, it’s probably a problem with your router configuration. Re-read the instructions above to see if you missed anything, search the web for help with “port forwarding”, and ask for help on [devChat](https://t.me/fuckgoogle).

## Configuring cyberd

Go to cyberd daemon folder, then go to `config` folder and open `config.toml` file for editing.

Find `peer to peer configuration options` section and edit `external_address` variable with your IP address and port `26656`

![peer_to_peer_config](https://ipfs.io/ipfs/QmQRqM4PbPt8cbDN49nAktT23XWfCixbfzfyUEkSyDUWYP)

Restart cyberd container.

---

We call to you, validators, with a proposal to forwarding port `26656` and make you validator-nodes available to the incoming connection.

Unfortunately, we can't provide all guides for port forwarding because of they different for each router. But if you faced on with some troubles feel free to contact us in our [devChat](https://t.me/fuckgoogle).
