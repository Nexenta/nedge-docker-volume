# nedge-docker-volume
Docker volume plugin for NexentaEdge

Building:
	% git clone https://github.com/Nexenta/nedge-docker-volume.git
	% make deps
	% make 

Running:
	% cp ndvol/daemon/ndvol.json /opt/nedge/etc/ccow/ndvol.json
	% bin/ndvol daemon start

