mount:
	mountpoint -q ./mnt || sudo mount /dev/nvme1n1p8 ./mnt
