# flightplan2litchimission
`fp2lm` is a command-line tool for converting the output generated by the [Flight Planner](https://github.com/JMG30/flight_planner) plugin for [QGIS](https://www.qgis.org/en/site/) to a [Litchi](https://flylitchi.com) mission.

Usage: `cat [FlightplannerMission].csv | fp2lm [-d] > [LitchiMission].csv`

`-d` sets the interval (distance) between projection centres in either meters or feet (feet are automatically converted to meters).  Accepted units are `m`, `ft`, `meters`, and `feet`.

For example, to set the interval between projection centres to 20 meters, use `fp2lm -d 20m`.

## Description

`fp2lm` reads a stream of waypoints generated by Flight Planner for QGIS and prints a stream of properly-structured Litchi Mission waypoints, line-by-line, to standard output (this is to allow subsequent stream editing using another tool, if desired).  I/O redirection may be used to capture the output in a new file for use by Litchi.

## Installing

Pre-compiled binaries are now available for macOS, Linux, and Windows 64-bit architectures (if you need something else, let me know) and are available in the project's `bin` directory.  Simply install the executable programme `fp2lm` (`fp2lm.exe` on Windows) in a suitable location and make sure its directory is included in your path.

### Mac OS

Open the Terminal, and copy the commands below.  Change any bracketed `[]` portions to reflect your particular environment.

1. Create a folder named `bin` in your home folder:


        mkdir ~/bin


2. Move the `fp2lm` binary from the download location to the newly-created `bin` folder:


        mv ~/Downloads/fp2lm ~/bin/fp2lm


3. Make the `fp2lm` programme executable:


        chmod u+x ~/bin/fp2lm


4. Open the Terminal and update your `PATH` to include the `bin` folder with the following command:


        export PATH=/Users/[your home folder]/bin:$PATH


5. To make the update to your `PATH` permanent, append the updated path to your user profile:


        echo "export PATH=/Users/[your home folder]/bin:$PATH" >> ~/.zshrc


6. You may now run `fp2lm` from the command line as described above at 'Usage:'.  For example, assuming you saved your QGIS Flightplanner flight plan as `FlightplannerMission.csv` on your desktop, and have determined you want twenty-meters between projection centres, you may run the following command:


        cat ~/Desktop/FlightplannerMission.csv | fp2lm -d 20m > ~/Desktop/LitchiMission.csv


Doing so will create a new file named `LitchiMission.csv` on your desktop, with the distance between projection centres set to twenty-meters, that may be uploaded to Litchi Mission Hub.

### Linux

Open the Terminal, and copy the commands below.  Change any bracketed `[]` portions to reflect your particular environment.

1. Create a folder named `bin` in your home folder:


        mkdir ~/bin


2. Move the `fp2lm` binary from the download location to the newly-created `bin` folder:


        mv ~/Downloads/fp2lm ~/bin/fp2lm


3. Make the `fp2lm` programme executable:


        chmod u+x ~/bin/fp2lm


4. Open the Terminal and update your `PATH` to include the `bin` folder with the following command:


        export PATH=/home/[your home folder]/bin:$PATH


5. To make the update to your `PATH` permanent, append the updated path to your user profile:


        echo "export PATH=/home/[your home folder]/bin:$PATH" >> ~/.profile


6. You may now run `fp2lm` from the command line as described above at 'Usage:'.  For example, assuming you saved your QGIS Flightplanner flight plan as `FlightplannerMission.csv` on your desktop, and have determined you want twenty-meters between projection centres, you may run the following command:


        cat ~/Desktop/FlightplannerMission.csv | fp2lm -d 20m > ~/Desktop/LitchiMission.csv


Doing so will create a new file named `LitchiMission.csv` on your desktop, with the distance between projection centres set to twenty-meters, that may be uploaded to Litchi Mission Hub.

### Windows

#### Prerequisites
- Ensure you are using Windows PowerShell, which is Unix/POSIX-compatible, for compatibility with certain command-line operations.

#### Step-by-Step Guide

1. **Create a Bin Directory**
   - Open Windows PowerShell.
   - Create a folder named `bin` (or name of your choice) in your home directory by executing:
     ```powershell
     mkdir $HOME\bin
     ```

2. **Move the fp2lm.exe Binary**
   - Move the `fp2lm.exe` binary from your download location to the newly-created `bin` folder. Assuming it is in your Downloads folder, use:
     ```powershell
     Move-Item $HOME\Downloads\fp2lm.exe $HOME\bin\fp2lm.exe
     ```

3. **Update the PATH Environment Variable**
   - Add the `bin` directory to your system's PATH environment variable. This can be done temporarily (just for the current session) by executing:
     ```powershell
     $Env:PATH += ";$HOME\bin"
     ```
   - For a permanent change, you will need to add `$HOME\bin` to the PATH environment variable through System Properties or by using the [System Environment Variables](https://docs.microsoft.com/en-us/windows/deployment/usmt/usmt-recognized-environment-variables) settings.

4. **Running fp2lm**
   - Now you can run `fp2lm` from the command line as described in the usage instructions. For example:
     ```powershell
     cat $HOME\Desktop\FlightplannerMission.csv | fp2lm -d 20m > $HOME\Desktop\LitchiMission.csv
     ```
        Alternatively, if you skipped step 3, specify the location of the fp2lm executable explicitly:
     
     ```powershell
     cat $HOME\Desktop\FlightplannerMission.csv | $HOME\bin\fp2lm -d 20m > $HOME\Desktop\LitchiMission.csv
     ```
   - This command assumes you have a file named `FlightplannerMission.csv` on your desktop. It will create a new file named `LitchiMission.csv` on your desktop, with the distance between projection centers set to twenty meters.

#### Notes
- The `cat` command is used to read the contents of the CSV file and is available by default in recent versions of PowerShell.
- If you encounter any issues, ensure that PowerShell is recognizing the `$HOME` variable correctly and that `fp2lm.exe` is located in the specified `bin` directory.




## Building from source

`go build fp2lm.go`

## Steps to produce a Litchi Mission using QGIS

This guide assumes the reader is already familiar with Litchi, but may need help with the workflow in QGIS.

1. Install the flight_planner plugin from QGIS → [Plugins] → [Manage and Install Plugins...] and search for ‘Flight Planner’.
2. Load the map layer of your choice.  To use Google Earth or OpenStreetMap, select ‘XYZ Tiles’ in your project's browser and add it as a layer to your project by double-clicking or right-clicking and selecting ‘Add Layer to Project’.
3. Scribe your Area of Interest (AoI) by creating a new shapefile layer from [Layer] → [Create Layer] → [New Shapefile Layer].  Select ‘Polygon’ as the Geometry type.  Select the desired points on the map. ℹ️ Depending on the CRS you are using, you may need to change the CRS of the AoI to work with Flight Planner — which requires measurements in meters.
4. Follow the [instructions](https://github.com/JMG30/flight_planner/wiki/Guide) for Flight Planner to plan your flight.  ℹ️ If you are using a DJI drone, you will probably need to add your own camera lens.  Consult the manufacture's specifications.
5. You will need to create latitude and longitude coordinates for use by Litchi.  Fortunately, QGIS makes this easy. With the flight plan generated, select the `waypoints` layer in the newly-created `flight_design` layer group.  Select [Vector] → [Geometry Tools] → [Add Geometry Attributes].  Select your AoI layer, and calculate the latitude and longitude coordinates using an appropriate CRS (for example, EPSG:4326).  Add the new layer to your project.  ℹ️ The newly-created layer will have two new fields for latitude and longitude called `xcoord` and `ycoord`.  You may verify the new values by right-clicking on the layer and selecting Open Attribute Table.
6. Export the new layer with latitude and longitude points added to a CSV file. ℹ️ If the steps were correctly followed, the exported file should have the following header: `️Waypoint Number,X [m],Y [m],Alt. ASL [m],Alt. AGL [m],xcoord,ycoord`
7. Measure the distance between projection centers (in the flight_design layer), you will supply this value to `fp2lm` in the final step.
8. Run `fp2lm` against the CSV file as described above with the distance between projection centres obtained in the step above set using the `-d` option.

**_NOTE:_** `fp2lm` expects CSV input in the form of navigation waypoints.  At the time of this writing, Litchi missions are limited to 99 waypoints, thus if waypoints are used to trigger photographs, or other actions, the allowed waypoints will be quickly consumed; therefore, for this workflow, waypoints are only used for course changes required to fly the grid.  Photograph intervals may be set expediently by measuring the distance between projection centres in QGIS, setting that distance in `fp2lm` using the `-d` flag, and configuring Litchi to photograph at equal distance intervals.  This is a work-around until Litchi adds support for more waypoints, but works very well.
