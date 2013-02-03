#!/usr/bin/env python

"""
Script that kills all Google Chrome processes (in Ubuntu).

This might be necessary for a variety of reasons. For example,
sometimes I have an issue where my synced profile can't be accessed.
The fix for this is to go into Chrome settings and reset the app password,
but the result won't take effect until you've killed the root Chrome
process, so you end up having to kill them all until you got it.
"""

import os
import re
import signal


UBUNTU_CHROME_PROCESS_BIN = '/opt/google/chrome/chrome'

def main():
    pids= [pid for pid in os.listdir('/proc') if pid.isdigit()]

    for pid in pids:
        name = open(os.path.join('/proc', pid, 'cmdline'), 'rb').read(100)
        if re.match(UBUNTU_CHROME_PROCESS_BIN, name):
            os.kill(int(pid), signal.SIGKILL)


if __name__ == '__main__':
    main()
