package utility

import (
	"math"
	"sort"
)

type ProvinceDistance struct {
	Province string
	Distance float64
}

var (
	provinces = map[string]struct {
		Latitude  float64
		Longitude float64
	}{
		"กรุงเทพมหานคร":   {13.7278956, 100.52412349999997},
		"กระบี่":          {8.0862997, 98.90628349999997},
		"กาญจนบุรี":       {14.0227797, 99.53281149999998},
		"กาฬสินธุ์":       {16.4314078, 103.5058755},
		"กำแพงเพชร":       {16.4827798, 99.52266179999992},
		"ขอนแก่น":         {16.4419355, 102.8359921},
		"จันทบุรี":        {12.61134, 102.10385459999998},
		"ฉะเชิงเทรา":      {13.6904194, 101.07795959999999},
		"ชลบุรี":          {13.3611431, 100.98467170000004},
		"ชัยนาท":          {15.1851971, 100.12512500000003},
		"ชัยภูมิ":         {15.8068173, 102.03150270000003},
		"ชุมพร":           {10.4930496, 99.18001989999993},
		"เชียงราย":        {19.9071656, 99.83095500000002},
		"เชียงใหม่":       {18.7877477, 98.99313110000003},
		"ตรัง":            {7.5593851, 99.61100650000003},
		"ตราด":            {12.2427563, 102.51747339999997},
		"ตาก":             {16.8839901, 99.12584979999997},
		"นครนายก":         {14.2069466, 101.21305110000003},
		"นครปฐม":          {13.8199206, 100.06216760000007},
		"นครพนม":          {17.392039, 104.76955079999993},
		"นครราชสีมา":      {14.9798997, 102.09776929999998},
		"นครศรีธรรมราช":   {8.4303975, 99.96312190000003},
		"นครสวรรค์":       {15.6930072, 100.12255949999997},
		"นนทบุรี":         {13.8621125, 100.51435279999998},
		"นราธิวาส":        {6.4254607, 101.82531429999995},
		"น่าน":            {18.7756318, 100.77304170000002},
		"บุรีรัมย์":       {14.9930017, 103.10291910000001},
		"ปทุมธานี":        {14.0208391, 100.52502759999993},
		"ประจวบคีรีขันธ์": {11.812367, 99.79732709999996},
		"ปราจีนบุรี":      {14.0509704, 101.37274389999993},
		"ปัตตานี":         {6.869484399999999, 101.25048259999994},
		"พระนครศรีอยุธยา": {14.3532128, 100.56895989999998},
		"พะเยา":           {19.1664789, 99.9019419},
		"พังงา":           {8.4407456, 98.51930319999997},
		"พัทลุง":          {7.6166823, 100.07402309999998},
		"พิจิตร":          {16.4429516, 100.34823289999997},
		"พิษณุโลก":        {16.8298048, 100.26149150000003},
		"เพชรบุรี":        {13.1111601, 99.93913069999996},
		"เพชรบูรณ์":       {16.4189807, 101.15509259999999},
		"แพร่":            {18.1445774, 100.14028310000003},
		"ภูเก็ต":          {7.9810496, 98.36388239999997},
		"มหาสารคาม":       {16.1850896, 103.30264609999995},
		"มุกดาหาร":        {16.542443, 104.72091509999996},
		"แม่ฮ่องสอน":      {19.2990643, 97.96562259999996},
		"ยโสธร":           {15.792641, 104.14528270000005},
		"ยะลา":            {6.541147, 101.28039469999999},
		"ร้อยเอ็ด":        {16.0538196, 103.65200359999994},
		"ระนอง":           {9.9528702, 98.60846409999999},
		"ระยอง":           {12.6833115, 101.23742949999996},
		"ราชบุรี":         {13.5282893, 99.81342110000003},
		"ลพบุรี":          {14.7995081, 100.65337060000002},
		"ลำปาง":           {18.2888404, 99.49087399999996},
		"ลำพูน":           {18.5744606, 99.0087221},
		"เลย":             {17.4860232, 101.72230020000006},
		"ศรีสะเกษ":        {15.1186009, 104.32200949999992},
		"สกลนคร":          {17.1545995, 104.1348365},
		"สงขลา":           {7.1756004, 100.61434699999995},
		"สตูล":            {6.6238158, 100.06737440000006},
		"สมุทรปราการ":     {13.5990961, 100.59983190000003},
		"สมุทรสงคราม":     {13.4098217, 100.00226450000002},
		"สมุทรสาคร":       {13.5475216, 100.27439559999993},
		"สระแก้ว":         {13.824038, 102.0645839},
		"สระบุรี":         {14.5289154, 100.91014210000003},
		"สิงห์บุรี":       {14.8936253, 100.39673140000002},
		"สุโขทัย":         {17.0055573, 99.82637120000004},
		"สุพรรณบุรี":      {14.4744892, 100.11771279999994},
		"สุราษฎร์ธานี":    {9.1382389, 99.32174829999997},
		"สุรินทร์":        {14.882905, 103.49371070000007},
		"หนองคาย":         {17.8782803, 102.74126380000007},
		"หนองบัวลำภู":     {17.2218247, 102.42603680000002},
		"อ่างทอง":         {14.5896054, 100.45505200000002},
		"อำนาจเจริญ":      {15.8656783, 104.62577740000006},
		"อุดรธานี":        {17.4138413, 102.78723250000007},
		"อุตรดิตถ์":       {17.6200886, 100.09929420000003},
		"อุทัยธานี":       {15.3835001, 100.02455269999996},
		"อุบลราชธานี":     {15.2286861, 104.85642170000006},
		"บึงกาฬ":          {18.3609104, 103.64644629999998},
	}
)

func DistanceCal(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515

	dist = dist * 1.609344

	return dist
}

func ProvinceCal(lat1 float64, lng1 float64) (string, string) {
	var distanceList []ProvinceDistance

	for province, coords := range provinces {
		distance := DistanceCal(lat1, lng1, coords.Latitude, coords.Longitude)
		distanceList = append(distanceList, ProvinceDistance{Province: province, Distance: distance})
	}

	sort.Slice(distanceList, func(i, j int) bool {
		return distanceList[i].Distance < distanceList[j].Distance
	})

	return distanceList[0].Province, distanceList[1].Province
}