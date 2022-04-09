#!/bin/python3
import shutil
from os import listdir,path,makedirs
import datetime
import exifread
from joblib import Parallel, delayed
import multiprocessing

# source = "E:/DCIM/100MSDCF/"
source = "C:/Users/danie/Pictures/tmp/"
dest = "C:/Users/danie/Pictures/dest/"
filelist = listdir(source)

num_cores = multiprocessing.cpu_count()
     
def copyphotos(source, dest, x):
  with open(source+x, "rb") as fh:
    tags = exifread.process_file(fh, stop_tag="EXIF DateTimeOriginal")
    creation_time = tags["EXIF DateTimeOriginal"]
    fh.close

  if creation_time != None:
    date = datetime.datetime.strptime(str(creation_time),"%Y:%m:%d %H:%M:%S")
    extension = path.splitext(source+x)[1]
    destpath = dest+str(extension[1:])+"/"+str(date.year)+"/"+str(date.month)+"/"+str(date.day)+"/"
    if not path.exists(destpath+x):
      makedirs(path.dirname(destpath), exist_ok=True)
    shutil.copy(source+x, destpath+x )
    print("copied file "+source+x+" to "+destpath+x)


Parallel(n_jobs=num_cores)(delayed(copyphotos)(source,dest,x)for x in filelist)