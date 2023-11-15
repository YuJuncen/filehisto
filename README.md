# Filehisto

A simple tool that collects the relationship between modify time and file size in a folder. You may check whehter there are some cold files in your disk with it.

## Build

```bash
go build -o bin/filehisto
```

## Usage

Check a folder's modtime and its 

```console
$ bin/filehisto --path /tmp
statistics for /tmp, total 20287 files
     6854h31m6s ago-    6397h33m36s ago  1.64%      ▍                 91.32MB
    6397h33m36s ago-     5940h36m5s ago  0%         ▏                 
     5940h36m5s ago-    5483h38m35s ago  0%         ▏                 
    5483h38m35s ago-     5026h41m4s ago  0%         ▏                 
     5026h41m4s ago-    4569h43m34s ago  0%         ▏                 
    4569h43m34s ago-     4112h46m3s ago  1.96e-05%  ▏                 1.09kB
     4112h46m3s ago-    3655h48m33s ago  0%         ▏                 
    3655h48m33s ago-     3198h51m2s ago  0%         ▏                 
     3198h51m2s ago-    2741h53m32s ago  0%         ▏                 
    2741h53m32s ago-     2284h56m1s ago  0%         ▏                 
     2284h56m1s ago-    1827h58m31s ago  3.56%      ▊                 198.4MB
    1827h58m31s ago-      1371h1m0s ago  1.45%      ▍                 80.54MB
      1371h1m0s ago-      914h3m30s ago  7.31%      █▋                407.3MB
      914h3m30s ago-      457h5m59s ago  16%        ███▌              891.2MB
      457h5m59s ago-          8m28s ago  70%        ███████████████▏  3.903GB
```

## Acknowledgement 

The histrogram part is copied from [aybabtme/uniplot](https://github.com/aybabtme/uniplot), licensed by MIT. I made a little change to it to fit the requirement.
