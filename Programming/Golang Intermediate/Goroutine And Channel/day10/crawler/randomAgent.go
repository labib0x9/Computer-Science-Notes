package main

import (
	"fmt"
	"math/rand"
	"strings"
)

var firefoxVersions = []float32{
	65.0, 66.0, 67.0, 68.0, 69.0, 70.0, 71.0, 72.0, 73.0, 74.0, 75.0, 76.0,
	77.0, 78.0, 79.0, 80.0, 81.0, 82.0, 83.0, 84.0, 85.0, 86.0, 87.0, 88.0,
	89.0, 90.0, 91.0, 92.0, 93.0, 94.0, 95.0, 96.0, 97.0, 98.0, 99.0, 100.0,
	101.0, 102.0, 103.0, 104.0, 105.0, 106.0, 107.0, 108.0, 109.0, 110.0,
	111.0, 112.0, 113.0, 114.0, 115.0, 116.0, 117.0, 118.0, 119.0, 120.0,
	121.0, 122.0, 123.0, 124.0, 125.0, 126.0, 127.0, 128.0, 129.0, 130.0,
	131.0, 132.0, 133.0, 134.0, 135.0, 136.0, 137.0, 138.0,
}

var chromeVersions = []string{
	"102.0.5005.115",
	"103.0.5060.53",
	"103.0.5060.66",
	"103.0.5060.114",
	"103.0.5060.134",
	"104.0.5112.79",
	"104.0.5112.80",
	"104.0.5112.81",
	"104.0.5112.101",
	"104.0.5112.102",
	"105.0.5195.52",
	"105.0.5195.53",
	"105.0.5195.54",
	"105.0.5195.102",
	"105.0.5195.125",
	"105.0.5195.126",
	"105.0.5195.127",
	"106.0.5249.61",
	"106.0.5249.62",
	"106.0.5249.91",
	"106.0.5249.103",
	"106.0.5249.119",
	"107.0.5304.62",
	"107.0.5304.63",
	"107.0.5304.68",
	"107.0.5304.87",
	"107.0.5304.88",
	"107.0.5304.106",
	"107.0.5304.107",
	"107.0.5304.110",
	"107.0.5304.121",
	"107.0.5304.122",
	"108.0.5359.71",
	"108.0.5359.72",
	"108.0.5359.94",
	"108.0.5359.95",
	"108.0.5359.98",
	"108.0.5359.99",
	"108.0.5359.124",
	"108.0.5359.125",
	"109.0.5414.74",
	"109.0.5414.75",
	"109.0.5414.87",
	"109.0.5414.119",
	"109.0.5414.120",
	"110.0.5481.77",
	"110.0.5481.78",
	"110.0.5481.96",
	"110.0.5481.97",
	"110.0.5481.100",
	"110.0.5481.104",
	"110.0.5481.177",
	"110.0.5481.178",
	"109.0.5414.129",
	"111.0.5563.64",
	"111.0.5563.65",
	"111.0.5563.110",
	"111.0.5563.111",
	"111.0.5563.146",
	"111.0.5563.147",
	"112.0.5615.49",
	"112.0.5615.50",
	"112.0.5615.86",
	"112.0.5615.87",
	"112.0.5615.121",
	"112.0.5615.137",
	"112.0.5615.138",
	"112.0.5615.165",
	"113.0.5672.63",
	"113.0.5672.64",
	"113.0.5672.92",
	"113.0.5672.93",
}

var edgeVersions = []string{
	"79.0.0.0,79.0.309.65",
	"80.0.0.0,80.0.361.62",
	"81.0.0.0,81.0.416.64",
	"83.0.0.0,83.0.478.37",
	"84.0.0.0,84.0.522.52",
	"85.0.0.0,85.0.564.41",
	"86.0.0.0,86.0.622.38",
	"87.0.0.0,87.0.664.41",
	"88.0.0.0,88.0.705.50",
	"89.0.0.0,89.0.774.45",
	"90.0.0.0,90.0.818.41",
	"91.0.0.0,91.0.864.37",
	"92.0.0.0,92.0.902.55",
	"93.0.0.0,93.0.961.38",
	"94.0.0.0,94.0.992.31",
	"95.0.0.0,95.0.1020.30",
	"96.0.0.0,96.0.1054.34",
	"97.0.0.0,97.0.1072.55",
	"98.0.0.0,98.0.1108.43",
	"99.0.0.0,99.0.1150.36",
	"121.0.0.0,121.0.2277.83",
	"122.0.0.0,122.0.2365.92",
	"123.0.0.0,123.0.2420.81",
	"124.0.0.0,124.0.2478.51",
	"125.0.0.0,125.0.2535.67",
	"126.0.0.0,126.0.2592.81",
	"127.0.0.0,127.0.2643.60",
	"128.0.0.0,128.0.2739.42",
	"129.0.0.0,129.0.2792.52",
	"130.0.0.0,130.0.2849.5",
	"131.0.0.0,131.0.2896.3",
	"132.0.0.0,132.0.2950.4",
	"133.0.0.0,133.0.3065.92",
	"134.0.0.0,134.0.3124.66",
	"135.0.0.0,135.0.3179.98",
	"136.0.0.0,136.0.3240.50",
	"137.0.0.0,137.0.3269.1",
	"138.0.0.0,138.0.3300.0",
	"139.0.0.0,139.0.3350.0",
	"140.0.0.0,140.0.3400.0",
}

var operaVersions = []string{
	"110.0.5449.0,96.0.4640.0",
	"110.0.5464.2,96.0.4653.0",
	"110.0.5464.2,96.0.4660.0",
	"110.0.5481.30,96.0.4674.0",
	"110.0.5481.30,96.0.4691.0",
	"110.0.5481.30,96.0.4693.12",
	"110.0.5481.77,96.0.4693.16",
	"110.0.5481.100,96.0.4693.20",
	"110.0.5481.178,96.0.4693.31",
	"110.0.5481.178,96.0.4693.50",
	"110.0.5481.192,96.0.4693.80",
	"111.0.5532.2,97.0.4711.0",
	"111.0.5532.2,97.0.4704.0",
	"111.0.5532.2,97.0.4697.0",
	"111.0.5562.0,97.0.4718.0",
	"111.0.5563.19,97.0.4719.4",
	"111.0.5563.19,97.0.4719.11",
	"111.0.5563.41,97.0.4719.17",
	"111.0.5563.65,97.0.4719.26",
	"111.0.5563.65,97.0.4719.28",
	"111.0.5563.111,97.0.4719.43",
	"111.0.5563.147,97.0.4719.63",
	"111.0.5563.147,97.0.4719.83",
	"112.0.5596.2,98.0.4756.0",
	"112.0.5596.2,98.0.4746.0",
	"112.0.5615.20,98.0.4759.1",
	"112.0.5615.50,98.0.4759.3",
	"112.0.5615.87,98.0.4759.6",
	"112.0.5615.165,98.0.4759.15",
	"112.0.5615.165,98.0.4759.21",
	"112.0.5615.165,98.0.4759.39",
}

var osStrings = []string{
	"Macintosh; Intel Mac OS X 10_13",
	"Macintosh; Intel Mac OS X 10_13_1",
	"Macintosh; Intel Mac OS X 10_13_2",
	"Macintosh; Intel Mac OS X 10_13_3",
	"Macintosh; Intel Mac OS X 10_13_4",
	"Macintosh; Intel Mac OS X 10_13_5",
	"Macintosh; Intel Mac OS X 10_13_6",
	"Macintosh; Intel Mac OS X 10_14",
	"Macintosh; Intel Mac OS X 10_14_1",
	"Macintosh; Intel Mac OS X 10_14_2",
	"Macintosh; Intel Mac OS X 10_14_3",
	"Macintosh; Intel Mac OS X 10_14_4",
	"Macintosh; Intel Mac OS X 10_14_5",
	"Macintosh; Intel Mac OS X 10_14_6",
	"Macintosh; Intel Mac OS X 10_15",
	"Macintosh; Intel Mac OS X 10_15_1",
	"Macintosh; Intel Mac OS X 10_15_2",
	"Macintosh; Intel Mac OS X 10_15_3",
	"Macintosh; Intel Mac OS X 10_15_4",
	"Macintosh; Intel Mac OS X 10_15_5",
	"Macintosh; Intel Mac OS X 10_15_6",
	"Macintosh; Intel Mac OS X 10_15_7",
	"Macintosh; Intel Mac OS X 11_0",
	"Macintosh; Intel Mac OS X 11_0_1",
	"Macintosh; Intel Mac OS X 11_1",
	"Macintosh; Intel Mac OS X 11_2",
	"Macintosh; Intel Mac OS X 11_2_1",
	"Macintosh; Intel Mac OS X 11_2_2",
	"Macintosh; Intel Mac OS X 11_2_3",
	"Macintosh; Intel Mac OS X 11_3",
	"Macintosh; Intel Mac OS X 11_3_1",
	"Macintosh; Intel Mac OS X 11_4",
	"Macintosh; Intel Mac OS X 11_5",
	"Macintosh; Intel Mac OS X 11_5_1",
	"Macintosh; Intel Mac OS X 11_5_2",
	"Macintosh; Intel Mac OS X 11_6",
	"Macintosh; Intel Mac OS X 11_6_1",
	"Macintosh; Intel Mac OS X 11_6_2",
	"Macintosh; Intel Mac OS X 11_6_3",
	"Macintosh; Intel Mac OS X 11_6_4",
	"Macintosh; Intel Mac OS X 11_6_5",
	"Macintosh; Intel Mac OS X 11_6_6",
	"Macintosh; Intel Mac OS X 11_6_7",
	"Macintosh; Intel Mac OS X 11_6_8",
	"Macintosh; Intel Mac OS X 11_7",
	"Macintosh; Intel Mac OS X 11_7_1",
	"Macintosh; Intel Mac OS X 11_7_2",
	"Macintosh; Intel Mac OS X 11_7_3",
	"Macintosh; Intel Mac OS X 11_7_4",
	"Macintosh; Intel Mac OS X 11_7_5",
	"Macintosh; Intel Mac OS X 11_7_6",
	"Macintosh; Intel Mac OS X 12_0",
	"Macintosh; Intel Mac OS X 12_0_1",
	"Macintosh; Intel Mac OS X 12_1",
	"Macintosh; Intel Mac OS X 12_2",
	"Macintosh; Intel Mac OS X 12_2_1",
	"Macintosh; Intel Mac OS X 12_3",
	"Macintosh; Intel Mac OS X 12_3_1",
	"Macintosh; Intel Mac OS X 12_4",
	"Macintosh; Intel Mac OS X 12_5",
	"Macintosh; Intel Mac OS X 12_5_1",
	"Macintosh; Intel Mac OS X 12_6",
	"Macintosh; Intel Mac OS X 12_6_1",
	"Macintosh; Intel Mac OS X 12_6_2",
	"Macintosh; Intel Mac OS X 12_6_3",
	"Macintosh; Intel Mac OS X 12_6_4",
	"Macintosh; Intel Mac OS X 12_6_5",
	"Macintosh; Intel Mac OS X 13_0",
	"Macintosh; Intel Mac OS X 13_0_1",
	"Macintosh; Intel Mac OS X 13_1",
	"Macintosh; Intel Mac OS X 13_2",
	"Macintosh; Intel Mac OS X 13_2_1",
	"Macintosh; Intel Mac OS X 13_3",
	"Macintosh; Intel Mac OS X 13_3_1",
	"Windows NT 10.0; Win64; x64",
	"Windows NT 5.1",
	"Windows NT 6.1; WOW64",
	"Windows NT 6.1; Win64; x64",
	"X11; Linux x86_64",
}

func genFirefoxUA() string {
	version := firefoxVersions[rand.Intn(len(firefoxVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%.1f) Gecko/20100101 Firefox/%.1f", os, version, version)
}

func genChromeUA() string {
	version := chromeVersions[rand.Intn(len(chromeVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", os, version)
}

func genEdgeUA() string {
	version := edgeVersions[rand.Intn(len(edgeVersions))]
	chromeVersion := strings.Split(version, ",")[0]
	edgeVersion := strings.Split(version, ",")[1]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 Edg/%s", os, chromeVersion, edgeVersion)
}

func genOperaUA() string {
	version := operaVersions[rand.Intn(len(operaVersions))]
	chromeVersion := strings.Split(version, ",")[0]
	operaVersion := strings.Split(version, ",")[1]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 OPR/%s", os, chromeVersion, operaVersion)
}

func userAgent() string {
	typee := rand.Intn(12324) % 4
	switch typee {
	case 0:
		return genFirefoxUA()
	case 1:
		return genChromeUA()
	case 2:
		return genEdgeUA()
	case 3:
		return genOperaUA()
	default:
		return genChromeUA()
	}
}
