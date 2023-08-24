yoloup
------------
yoloup is a simple commandline utility to update yolo annotation files based on a new class file.

### Why should I use yoloup?
Many offline labelling does not provide a way to save your annotation projects in between. Every time you start a new session to label a new set of images, the class index will start from 0 rather than your last class index + 1. That means, a new class file will be generated, which has no relationship with your last class file. For example, "elephant" indexed at 0 in your last class file might end up indexed at 5 in your current class file.

### Usage
```{sh}
Usage: yoloup [OPTIONS] [orginal_class_file] [updated_class_file] [yolo_annotation_file ...]
  -h    Print help information.
  -i    Request confirmation before attempting to update each file
```

### Examples
Let's say there is a yolo annotation file "RCNX0001.txt"
```
0 0.449951 0.490972 0.232910 0.276389
1 0.348754 0.402843 0.252323 0.306251
```
, and the original class file "labels.txt"
```
elephant
lion
horse
```
Suppose you now have an updated class file "target.txt", which include more classes than labels.txt but the order is messed up. For example, elephant is move from index 0 to index 2.
```
donkey
tiger
elephant
dog
cat
lion
horse
```
And you want use a new set of indexes based on a new class file "target.txt" instead of "labels.txt" to include more classes, you can do the following:
```
yoloup ./labels.txt ./target.txt ./RCNX0001.txt
```
This will change RCNX0001.txt to:
```
2 0.449951 0.490972 0.232910 0.276389
5 0.348754 0.402843 0.252323 0.306251
```





