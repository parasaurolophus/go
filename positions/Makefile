plot: clean build
	./positions > positions.csv
	gnuplot -p -e "D='`date`'" positions.gnuplot

clean:
	go clean
	if [ -f positions.csv ]; then rm positions.csv; fi

build:
	go build
