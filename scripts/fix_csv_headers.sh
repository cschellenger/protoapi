#!/bin/bash

if [ -f Crime_Data_from_2020_to_Present.csv ]; then
    cat proto/la_crime_headers.csv > crime.csv
    tail -n +2 Crime_Data_from_2020_to_Present.csv >> crime.csv
    echo "Wrote 'crime.csv'"
else
    echo "Missing 'Crime_Data_from_2020_to_Present.csv'"
    exit 1
fi

