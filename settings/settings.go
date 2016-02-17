package settings

const URL string = "http://eikaiwa.dmm.com/"

type Teacher struct {
	Id      int
	Name    string
	Country string
}

//2015-01-08 update
var TEACHERS_ID []Teacher = []Teacher{
	{Id: 6214, Name: "Aleksandra S", Country: "Serbia"},
	{6466, "Rachel L", "Shingapore"},
	{3293, "Edina", "Serbia"},
	{4107, "Milica Ml", "Serbia"},
	{8519, "Marine", "France"},
	{8261, "Ela T", "Germany"},
	{453, "Meryem", "America"},
	{4287, "Sana", "America"},
	//{5200, "Meisha", "America"},
	{5416, "Alyssa J", "America"},
	{5417, "Ashley G", "America"},
	{5482, "N Mika", "America"},
	{5922, "Mary Kate", "America"},
	{6051, "Michaela K", "America"},
	{6542, "Keira", "America"},
	{6665, "Claire K", "America"},
	{6926, "Katie S", "America"},
	{7355, "Mei Li", "America"},
	{7545, "Amanda Y", "America"},
	{7582, "Kittay", "America"},
	{8099, "Dilia", "America"},
	{8475, "Lindsey", "America"},
	{8542, "Miranda Faye", "America"},
	{8562, "Carla M", "America"},
	{8602, "Abbey Claire", "America"},
	{8703, "Mary C", "America"},
	//{6281, "Monica U", "Canada"},
	{398, "Jules", "UK"},
	{1670, "Jessica", "UK"},
	{425, "Mel", "Australia"},
}
