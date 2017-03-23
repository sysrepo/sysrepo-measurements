#!/usr/bin/python3
#Importing the modules
#pip install configparser
import os
import ConfigParser
import time
import json
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import numpy as np

config = ConfigParser.ConfigParser()
config.read("./config.txt")

for section in config.sections():

    fig = plt.figure()
    ax = fig.add_subplot(111)

    ax.plot(json.loads(config.get(section,"x")), json.loads(config.get(section, "y")), config.get(section, "color"), label=config.get(section, "label"))

    # Now add the legend.
    ax.legend(loc='upper left')

    fig.savefig(section + '.png')


