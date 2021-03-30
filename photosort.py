#!/bin/python3
import shutil
from PIL import Image
from os import listdir,path,makedirs
import datetime
import PySimpleGUI as sg

sourcepath = "./source/"
sourcelist = listdir(sourcepath)

sg.theme('DarkTeal2')   # Add a touch of color
# All the stuff inside your window.
layout = [  [sg.Text('Some text on Row 1')],
            [sg.Text('Enter something on Row 2'), sg.InputText()],
            [sg.Button('Ok'), sg.Button('Cancel')] ]

# Create the Window
window = sg.Window('Window Title', layout)
# Event Loop to process "events" and get the "values" of the inputs
while True:
    event, values = window.read()
    if event == sg.WIN_CLOSED or event == 'Cancel': # if user closes window or clicks cancel
        break
    print('You entered ', values[0])

window.close()




def copyphotos(sourcepath, destpath):
  for x in sourcelist:
    im = Image.open(sourcepath+x)
    exif = im.getexif()
    creation_time = exif.get(36867)
    
    if creation_time != None:
      date = datetime.datetime.strptime(creation_time,"%Y:%m:%d %H:%M:%S")
      extension = path.splitext(sourcepath+x)[1]
      destpath = "./dest/"+str(date.year)+"/"+str(date.month)+"/"+str(date.day)+"/"
  
      if not path.exists(destpath+x):
        makedirs(path.dirname(destpath), exist_ok=True)
        shutil.copy(sourcepath+x, destpath+x )
        print("copied file "+sourcepath+x+" to "+destpath+x)
