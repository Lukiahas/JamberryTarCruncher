package main

import (
	"log"
	"os"
	"io/ioutil"
	"encoding/csv"
	"fmt"
	"strconv"
	"regexp"
)

var ConsByLevel = [][]Consultant{}
var RankMap = make(map[string]Rank)
var RankNames = []string{}
var TMMap = make(map[string]Consultant)
var ThisConsultant Consultant
var ThisRank Rank
var ForceRank = int64(-1)

func main() {

	// Init the RankMap with predefined Ranks. 
	GenerateRanks()

	var file string
	
	if len(os.Args) == 2{
		file = os.Args[1];
	}else if len(os.Args) == 3{
		file = os.Args[1]
		ForceRank, _ = strconv.ParseInt(os.Args[2], 10, 0)
		if ForceRank > 13 || ForceRank < 1{
			log.Fatal("Please provide a valid second argument (1 <= x <= 13) where 1 is Consultant and 13 is Elite Executive")
		}
		//  = int(i)
		
	}else{
		//file = "Export.csv"
		
		// Look for latest .csv file in the current directory. 
		files, err := ioutil.ReadDir(".")
		if err != nil { log.Fatal(err) }
		
		re := regexp.MustCompile("\\.csv$")
		
		var newest os.FileInfo
		
		
		for _, file := range files{
			if re.MatchString(file.Name()){
				if newest == nil || file.ModTime().After(newest.ModTime()){
					newest = file
				}
			}
		}
		
		if newest == nil{
			log.Fatal("No TAR file found. Please put your latest TAR .csv file in the folder you are executing from.")
		}else{
			file = newest.Name()
			fmt.Print("File being read:\t", file, "\n")
		}
		
	}
	
	// Read CSV File (AKA the TAR)
	ReadCSV(file)
	
	// Load the global variables with the Consultant to do the calculations for.
	ThisConsultant = ConsByLevel[0][0]
	if ForceRank > 0{
		ThisRank = RankMap[RankNames[ForceRank - 1]]
	}else{
		ThisRank = RankMap[ThisConsultant.PayRank.Title]
	}
	
	fmt.Print("Full Name:\t\t", ThisConsultant.FullName(), "\n")
	fmt.Print("Actual Rank:\t\t", ThisConsultant.PayRank.Title, "\n")
	fmt.Print("Calculated Rank:\t", ThisRank.Title, "\n\n")

	// Do the calculations
	pb := PersonalBonus()
	lb := LevelBonus()
	gb := GenerationBonus()
	fb := FastStartBonus()
	
	// Print the numbers

	fmt.Printf("Personal Bonus:\t\t$%.2f\nLevel Bonus:\t\t$%.2f\nGeneration Bonus:\t$%.2f\nFast Start Bonus:\t$%.2f\n\nTotal Bonus Check:\t$%.2f\n\n",pb, lb, gb, fb,  pb + lb + gb + fb)
	
	// Pause before exiting to let the user view the numbers.
	fmt.Printf("Press Enter to exit.\n\n")
	var temp string
	_,_ = fmt.Scanln(&temp)
}

// Generate the Consultant ranks.
func GenerateRanks(){ //RankMap *map[string]Rank
	RankMap["Consultant"] = NewRank("Consultant","C",0.03,0.00,0.00,0.05,0.00,0.00,0.00,0.00,0.00,0.00,0.00,1)
	RankMap["Advanced Consultant"] = NewRank("Advanced Consultant","AC",0.05,0.00,0.00,0.05,0.05,0.00,0.00,0.00,0.00,0.00,0.00,2)
	RankMap["Senior Consultant"] = NewRank("Senior Consultant","SC",0.07,0.03,0.00,0.05,0.05,0.00,0.00,0.00,0.00,0.00,0.00,3)
	RankMap["Lead Consultant"] = NewRank("Lead Consultant","LC",0.10,0.05,0.00,0.05,0.05,0.00,0.00,0.00,0.00,0.00,0.00,4)
	RankMap["Senior Lead Consultant"] = NewRank("Senior Lead Consultant","SLC",0.12,0.06,0.03,0.05,0.05,0.00,0.00,0.00,0.00,0.00,0.00,5)
	RankMap["Premier Consultant"] = NewRank("Premier Consultant","PC",0.12,0.07,0.04,0.05,0.05,0.00,0.00,0.00,0.00,0.00,0.00,6)
	RankMap["Team Manager"] = NewRank("Team Manager","TM",0.12,0.07,0.05,0.05,0.05,0.02,0.00,0.00,0.00,0.00,0.00,7)
	RankMap["Senior Team Manager"] = NewRank("Senior Team Manager","STM",0.12,0.07,0.05,0.05,0.05,0.02,0.03,0.00,0.00,0.00,0.00,8)
	RankMap["Executive"] = NewRank("Executive","E",0.12,0.07,0.05,0.05,0.05,0.02,0.04,0.03,0.00,0.00,0.00,9)
	RankMap["Senior Executive"] = NewRank("Senior Executive","SE",0.12,0.07,0.05,0.05,0.05,0.02,0.04,0.04,0.03,0.00,0.00,10)
	RankMap["Lead Executive"] = NewRank("Lead Executive","LE",0.12,0.07,0.05,0.05,0.05,0.02,0.04,0.04,0.04,0.03,0.00,11)
	RankMap["Premier Executive"] = NewRank("Premier Executive","PE",0.12,0.07,0.05,0.05,0.05,0.02,0.04,0.04,0.04,0.04,0.03,12)
	RankMap["Elite Executive"] = NewRank("Elite Executive","EE",0.12,0.07,0.05,0.05,0.05,0.02,0.04,0.04,0.04,0.04,0.04,13)
	
	RankNames = []string{"Consultant",
		"Advanced Consultant",
		"Senior Consultant",
		"Lead Consultant",
		"Senior Lead Consultant",
		"Premier Consultant",
		"Team Manager",
		"Senior Team Manager",
		"Executive",
		"Senior Executive",
		"Lead Executive",
		"Premier Executive",
		"Elite Executive"}
	
}

// Read the CSV (AKA the TAR)
func ReadCSV(file string){
	csvfile, err := os.Open(file)
	if err != nil { log.Fatal(err) }

	defer csvfile.Close()
	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = 0 // see the Reader struct information below Sets "fields per record" to the number of fields in the first record.
	reader.LazyQuotes = true
	
	var consultant Consultant
	each, err := reader.Read()
	if len(each) != 29 || each[5] != "Phone" { log.Fatal("Wrong (old) format or wrong csv width") } //Check first row. Should be titles of fields.
	for each, err = reader.Read() ; err == nil; each, err = reader.Read() {
		//fmt.Println( each[0]) //Debug
		consultant.Number = ParseInt(each[0])
		consultant.DownlineLevel  = ParseInt(each[1])
		consultant.FirstName = each[2]
		consultant.LastName = each[3]
		consultant.Email = each[4]
		consultant.Phone = each[5]
		consultant.StartDate = each[6]
		consultant.Status = CleanHtmlTags(each[7])
		consultant.ActiveLegs = ParseInt(each[8])
		consultant.HighestLegRank = NewRankTitle(each[9])
		consultant.PayRank = NewRankTitle(each[10])
		consultant.CareerTitle = NewRankTitle(each[11])
		consultant.PRV  = ParseFloat(each[12])
		consultant.CV  = ParseFloat(each[13])
		consultant.TRV  = ParseFloat(each[14])
		consultant.DRV  = ParseFloat(each[15])
		consultant.UplineTM = each[16]
		consultant.Address = each[17]
		consultant.City = each[18]
		consultant.State = each[19]
		consultant.Zip  = each[20]
		consultant.Country = each[21]
		consultant.SponsoredThisMonth  = ParseInt(each[22])
		consultant.NumberInDownline  = ParseInt(each[23])
		consultant.LastLogin = each[24]
		consultant.Sponsor = each[25]
		consultant.SponsorEmail = each[26]
		consultant.Type = each[27]
		
		// Put it in the 2d array. First D is downline level. Second D is each consultant in that level.
		if len(ConsByLevel) == consultant.DownlineLevel {
			level := []Consultant{consultant}
			ConsByLevel = append(ConsByLevel, level)
		}else if len(ConsByLevel) > consultant.DownlineLevel {
			ConsByLevel[consultant.DownlineLevel] = append(ConsByLevel[consultant.DownlineLevel], consultant)
		}
		
		// Fill the Manager Map and calculate Generation for Generation Override Bonus. 
		if consultant.IsManager(){
			if val, ok := TMMap[consultant.UplineTM]; ok {
				consultant.Generation = val.Generation + 1				
			} else {
				consultant.Generation = 0
			}
			TMMap[consultant.FullName()] = consultant
		}
	}
	//fmt.Println(consultant.Number) //debug
}

// Calculates Personal Bonus
func PersonalBonus() float32{
	var bonus float32
	if ThisConsultant.PRV < 200{
		bonus = 0.0
	} else if ThisConsultant.PRV < 500{
	bonus = 0.03
	} else if ThisConsultant.PRV < 1500{
		bonus = 0.05
	} else if ThisConsultant.PRV < 3000{
		bonus = 0.08
	} else {
		bonus = 0.1
	}
	return float32(ThisConsultant.PRV) * bonus
}

// Calculates Fast Start Bonus
func FastStartBonus() float32{
	bonus := float32(0)
	for i := 1; i < 3; i++{
		for _, consult := range ConsByLevel[i]{
			if consult.Type == "Fast-Start"{
				bonus += float32(consult.CV) * ThisRank.FastStart[i]
			}
		}
	}
	return bonus
}

// Calculates Generation Override Bonus
func GenerationBonus() float32{
	bonus := float32(0)
	for _, tm := range TMMap{
		bonus += ThisRank.GenOverride[tm.Generation] * float32(tm.CV)
		for i := 0; i < len(ConsByLevel); i++{
			for _, consult := range ConsByLevel[i]{
				if !consult.IsManager() && tm.FullName() == consult.UplineTM{
					bonus += ThisRank.GenOverride[tm.Generation] * float32(consult.CV)
				}
			}
		}
		// log.Print("DownlineCount for: ", tm.FullName(), " is: ", downlinecount) //Debug
	}
	return bonus
}

// Calculates Level Bonus.
func LevelBonus() float32{
	var multiplier float32
	var maxLevel int
	totalBonus := float32(0)
	if len(ConsByLevel) > 4{
		maxLevel = 4
	}else{
		maxLevel = len(ConsByLevel)
	}
	for i := 1; i < maxLevel; i++{
		multiplier = ThisRank.Override[i]
		for _, downline := range ConsByLevel[i]{
			totalBonus += float32(downline.CV) * multiplier
		}
	}
	return totalBonus
}

// Parse ints from the TAR.
func ParseInt(input string) int{
	re := regexp.MustCompile("[^0-9.]")
	i, err := strconv.ParseInt(re.ReplaceAllString(input, ""), 10, 0)
	if err != nil { fmt.Println(i);fmt.Println(input); log.Fatal(err) }
	return int(i)
}

// Parse floats from the TAR.
func ParseFloat(input string) float32{
	re := regexp.MustCompile("[^0-9.]")
	i, err := strconv.ParseFloat(re.ReplaceAllString(input, ""), 0)
	if err != nil { fmt.Println(i);fmt.Println(input); log.Fatal(err) }
	return float32(i)
}

// Removes <nobr> tags from TAR String.
func CleanHtmlTags(input string) string{
	re := regexp.MustCompile("</?[a-z]*>")
	return re.ReplaceAllString(input, "")
}

type Consultant struct {
	Number int
	DownlineLevel int
	FirstName string
	LastName string
	Email string
	Phone string
	StartDate string
	Status string
	ActiveLegs int
	HighestLegRank RankTitle
	PayRank RankTitle
	CareerTitle RankTitle
	PRV float32
	CV float32
	TRV float32
	DRV float32
	UplineTM string
	Generation int //?
	Address string
	City string
	State string
	Zip string
	Country string
	SponsoredThisMonth int
	NumberInDownline int
	LastLogin string
	Sponsor string
	SponsorEmail string
	Type string
	//AttendingConference bool //Last field ignored.
}

//Check if Consultant is a TM or above. 
func (con *Consultant) IsManager() bool{
	if len(RankMap) == 0 { return false }
	return RankMap[con.CareerTitle.Title].Level >= RankMap["Team Manager"].Level
}

//Concat Full Name.
func (con *Consultant) FullName() string{
	return con.FirstName + " " + con.LastName
}

type RankTitle struct{
	Title string
	Number int
}

func NewRankTitle(title string) RankTitle{
	rankTitle := new(RankTitle)
	re := regexp.MustCompile("[0-9]{2}|([A-Z][a-zA-Z ]*[a-z])")
	splitTitle := re.FindAllString(title, -1)
	i, err := strconv.ParseInt(splitTitle[0], 10, 0)
	if err != nil { log.Fatal(err) }
	rankTitle.Number = int(i)
	rankTitle.Title = splitTitle[1]
	return *rankTitle
}


type Rank struct {
	Title string
	Abbreviation string
	Override map[int]float32
	FastStart map[int]float32
	GenOverride map[int]float32
	Level int
}

func NewRank(t, a string, o1, o2, o3, f1, f2, g0, g1, g2, g3, g4, g5 float32, level int) Rank{
	rank := new(Rank)
	rank.Title = t
	rank.Abbreviation = a
	rank.Override = make(map[int]float32)
	rank.Override[1] = o1
	rank.Override[2] = o2
	rank.Override[3] = o3
	rank.FastStart = make(map[int]float32)
	rank.FastStart[1] = f1
	rank.FastStart[2] = f2
	rank.GenOverride = make(map[int]float32)
	rank.GenOverride[0] = g0
	rank.GenOverride[1] = g1
	rank.GenOverride[2] = g2
	rank.GenOverride[3] = g3
	rank.GenOverride[4] = g4
	rank.GenOverride[5] = g5
	rank.Level = level
	return *rank
}