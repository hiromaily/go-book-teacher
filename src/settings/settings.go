package settings

const URL string = "http://eikaiwa.dmm.com/"

type Teacher struct {
	Id      int
	Name    string
	Country string
}

var TEACHERS_ID []Teacher = []Teacher{
	{Id: 6214, Name: "Aleksandra S", Country: "Serbia"},
	//{5453, "test teacher", "Shingapore"},
	{6466, "Rachel L", "Shingapore"},
	{5537, "Eleasha", "Malaysia"},
	{453, "Meryem", "America"},
	{4287, "Sana", "America"},
	{5200, "Meisha", "America"},
	{5416, "Alyssa J", "America"},
	{5417, "Ashley G", "America"},
	{5419, "Angela Rachel", "America"},
	{5482, "N Mika", "America"},
	{5817, "Rina P", "America"},
	{5922, "Mary Kate", "America"},
	{5961, "Deanna TN", "America"},
	{6051, "Michaela K", "America"},
	{6542, "Keira", "America"},
	{6665, "Claire K", "America"},
	{6736, "Jen", "America"},
	{6740, "Francesca", "America"},
	{6926, "Katie S", "America"},
	{7355, "Mei Li", "America"},
	{7448, "Mia", "America"},
	{7449, "Ketsia", "America"},
	{907, "Corrine", "UK"},
	{1670, "Jessica", "UK"},
	{4515, "Camila", "Canada"},
	{425, "Mel", "Australia"},
}
