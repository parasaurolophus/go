set datafile separator ','

TITLE = sprintf('Sun Positions %s', D)
set title TITLE offset char 0, char -1
set xdata time         # x axis is time data
set timefmt '%H:%M:%S' # input time string format
set format x '%H:%M'   # output time format
set key autotitle columnhead

plot 'positions.csv' using 1:2 with lines, '' using 1:3 with lines
