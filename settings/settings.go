package settings

const URL string = "http://eikaiwa.dmm.com/"

type Teacher struct {
	Id      int
	Name    string
	Country string
}

//2015-02-24 update
var TEACHERS_ID []Teacher = []Teacher{
	{Id: 6214, Name: "Aleksandra S", Country: "Serbia"},
	{1411, "Jekaterina", "Latvia"},
	{3293, "Edina", "Serbia"},
	{4107, "Milica Ml", "Serbia"},
	{4565, "Mana", "Serbia"},
	{4806, "Jekica", "Serbia"},
	{4808, "Joxyly", "Serbia"},
	{5656, "Lavinija", "Serbia"},
	{5809, "Sandra Z", "Serbia"},
	{6550, "Yovana", "Serbia"},
	{8160, "Gaja", "Serbia"},
	{7093, "Rita M", "Portugal"},
	{6466, "Rachel L", "Shingapore"},
	{8519, "Marine", "France"},
	//{8261, "Ela T", "Germany"},
	{453, "Meryem", "America"},
	{4287, "Sana", "America"},
	{5416, "Alyssa J", "America"},
	{5417, "Ashley G", "America"},
	{5482, "N Mika", "America"},
	{5922, "Mary Kate", "America"},
	{6051, "Michaela K", "America"},
	{6542, "Keira", "America"},
	{6665, "Claire K", "America"},
	{6926, "Katie S", "America"},
	{7238, "Aminah", "America"},
	{7355, "Mei Li", "America"},
	{7545, "Amanda Y", "America"},
	{8099, "Dilia", "America"},
	{8542, "Miranda Faye", "America"},
	{8562, "Carla M", "America"},
	{8602, "Abbey Claire", "America"},
	{8703, "Mary C", "America"},
	{8720, "Michelle Lynn", "America"},
	//{8810, "Kobi", "America"},
	//{8827, "Anne Alexis", "America"},
	{8869, "Kenda", "America"},
	//{8868, "Kimberly M", "Canada"},
	{398, "Jules", "UK"},
	//{1670, "Jessica", "UK"},
	{425, "Mel", "Australia"},
	{9016, "JessieRB", "New Zealand"},
}
