# Supported GPU list for pussy validators

In our `pussy protocol` implementation on `GO` proof of relevance root hash is computed on Cuda GPUs every round as the best way to calculate merkle tree faster. We need to load the whole graph in memory for calculating that's why memory volume is important. GPU with 6Gb memory can calculate graph with ~200 M links.

|GPU|Supported|Tested|CUDA cores|Memory|Year of production|
|---|---|---|---|---|---|
|[GEFORCE RTX 2080 Ti](https://www.nvidia.com/en-us/geforce/graphics-cards/rtx-2080-ti/)|:white_check_mark:|:x:|4352|11GB GDDR 6|2018|
|[GEFORCE RTX 2080](https://www.nvidia.com/en-us/geforce/graphics-cards/rtx-2080/)|:white_check_mark:|:x:|4352|11GB GDDR 6|2018|
|[GEFORCE RTX 2070](https://www.nvidia.com/en-us/geforce/graphics-cards/rtx-2070/)|:white_check_mark:|:x:|2304|8 GB GDDR6|2019|
|[GeForce RTX 2060](https://www.nvidia.com/en-us/geforce/graphics-cards/rtx-2060/)|:white_check_mark:|:x:|1920|6 GB GDDR6|2019|
|[GEFORCE GTX 1660 Ti](https://www.nvidia.com/en-us/geforce/graphics-cards/gtx-1660-ti/)|:white_check_mark:|:x:|1536|6GB GDDR6|2019|
|[GEFORCE GTX 1660](https://www.nvidia.com/en-us/geforce/graphics-cards/gtx-1660-ti/)|:white_check_mark:|:x:|1408|6GB GDDR5|2019|
|[GEFORCE GTX 1650](https://www.nvidia.com/en-us/geforce/graphics-cards/gtx-1650/)|:white_check_mark:|:white_check_mark:|896|4GB GDDR5|2019|
|[GeForce GTX 1080](https://www.nvidia.com/en-us/geforce/products/10series/geforce-gtx-1080/)|:white_check_mark:|:white_check_mark:|2560|8 GB GDDR5X|2016|
|[GeForce GTX 980](https://www.geforce.com/hardware/desktop-gpus/geforce-gtx-980/specifications)|:white_check_mark:|:x:|2048|4 GB GDDR5|2014|
|[TITAN Xp](https://www.nvidia.com/en-us/titan/titan-xp/)|:white_check_mark:|:x:|3840|12 GB GDDR5|2017|
|[GeForce GTX 1080 Ti](https://www.nvidia.com/en-us/geforce/products/10series/geforce-gtx-1080-ti/)|:white_check_mark:|:x:|3584|11 GB GDDR5X|2017|
|[GeForce GTX 980 Ti](https://www.nvidia.com/en-us/geforce/products/10series/geforce-gtx-1080-ti/)|:white_check_mark:|:x:|2816|6 GB GDDR5|2015|
|[GeForce GTX 1070 Ti](https://www.nvidia.com/en-us/geforce/products/10series/geforce-gtx-1070-ti/)|:white_check_mark:|:white_check_mark:|2432|8 GB GDDR5|2017|
|[GeForce GTX 1070](https://www.nvidia.com/en-us/geforce/products/10series/geforce-gtx-1070-ti/)|:white_check_mark:|:white_check_mark:|1920|8 GB GDDR5|2016|
|[GeForce GTX 970](https://www.nvidia.com/en-us/geforce/products/10series/geforce-gtx-1070-ti/)|:white_check_mark:|:x:|1664|4 GB GDDR5|2015|
|[GEFORCE GTX 1060 6GB](https://www.nvidia.com/en-us/geforce/products/10series/geforce-gtx-1060/)|:white_check_mark:|:white_check_mark:|1280|6 GB GDDR5|2016|
|[GeForce GTX 1050 Ti 4GB](https://www.nvidia.com/en-us/geforce/products/10series/geforce-gtx-1050/)|:white_check_mark:|:white_check_mark:|768|4 GB GDDR5|2016|
|[GeForce GTX 1050 4GB mobile](https://www.techpowerup.com/gpu-specs/asus-gtx-1050-mobile-4-gb.b5497)|:white_check_mark:|:white_check_mark:|640|4 GB GDDR5|2020|
|[GeForce GTX 745 (OEM) 4GB](https://www.geforce.com/hardware/desktop-gpus/geforce-gtx-745-oem/specifications)|:white_check_mark:|:x:|768|4 GB GDDR3|2014|
|[GeForce GTX TITAN X](https://www.geforce.com/hardware/desktop-gpus/geforce-gtx-titan-x/specifications)|:white_check_mark:|:x:|3072|12 GB GDDR5|2016|
|[GeForce GTX TITAN Z](https://www.geforce.com/hardware/desktop-gpus/geforce-gtx-titan-z/specifications)|:white_check_mark:|:x:|5760|12 GB GDDR5|2014|
|[GeForce GTX TITAN Black](https://www.geforce.com/hardware/desktop-gpus/geforce-gtx-titan-black/specifications)|:white_check_mark:|:x:|2880|6 GB GDDR5|2014|
|[GeForce GTX 770](https://www.geforce.com/hardware/desktop-gpus/geforce-gtx-770/specifications)|:white_check_mark:|:x:|1536|4 GB GDDR5|2013|
|[Nvidia Quadro K2200](https://www.nvidia.com/content/dam/en-zz/Solutions/design-visualization/documents/75509_DS_NV_Quadro_K2200_US_NV_HR.pdf)|:white_check_mark:|:white_check_mark:|640|4 GB GDDR5|2016|

If you have used some GPU from `column` supported but without :white_check_mark: at `tested` column please submit a pull request with corrections. If you have tested GPU and it's not contained in that list submit PR too.

**Note** If you using some old cards (like GTX 770, or older) make sure your card will be supported by al least **v.410** of NVIDIA diver for Linux.
