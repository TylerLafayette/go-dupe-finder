# go-dupe-finder
A small CLI application made to search a specific directory for duplicate files. This was made as an example program for an introductory Golang talk, and is not without bugs or perfectly designed. An example `images` folder is included with a few duplicates to try the program on. All the photos are from Unsplash and are labeled with the respective artist's name. Thanks to the artists of all these great pictures!

## The design
```
Input directory name
	-> get the names of all files in the directory
	-> divvy up the work between the maximum thread count
	-> each forked process will open each file and hash its contents
	-> once finished, check each key of the hash map for duplicates 
```

## How to use it
Simply run the command with your desired directory and maximum number of threads (Goroutines).  
For example,
```bash
./go-dupe-finder --path images --threads 100
```
