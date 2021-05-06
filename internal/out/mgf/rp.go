package mgf

var (
	RP = map[string]string{

		"b-102": `ENDCUT_TYPE=E213RS
START_OF_PARAMS
NO_OF_PARAMS=12
A=<%= data["r"] %>
B=<%= data["r"] %>
C=0
R1=<%= data["r"] %>
CD=12
V1=<%= 90.0 + data["angle"] %>
V2=30
V4=30
GAMMAO=30
HO=<%= t - 2.0 %>
GAMMAU=0
HU=0
END_OF_PARAMS`,

		"b-103": `ENDCUT_TYPE=E213RS
START_OF_PARAMS
NO_OF_PARAMS=12
A=<%= data["r"] %>
B=<%= data["r"] %>
C=0
R1=<%= data["r"] %>
CD=-1
V1=<%= 90.0 + data["angle"] %>
V2=30
V4=30
GAMMAO=30
HO=<%= t / 2.0 - 1.0 %>
GAMMAU=0
HU=0
END_OF_PARAMS`,

		"b-104": `ENDCUT_TYPE=E213RS
START_OF_PARAMS
NO_OF_PARAMS=12
A=<%= data["r"] %>
B=<%= data["r"] %>
C=0
R1=<%= data["r"] %>
CD=-1
V1=<%= 90.0 + data["angle"] %>
V2=45
V4=45
GAMMAO=45
HO=<%= t - 2.0 %>
GAMMAU=0
HU=0
END_OF_PARAMS`,

		"b-106": `ENDCUT_TYPE=E213RS
START_OF_PARAMS
NO_OF_PARAMS=12
A=<%= data["r"] %>
B=<%= data["r"] %>
C=0
R1=<%= data["r"] %>
CD=13
V1=<%= 90.0 + data["angle"] %>
V2=30
V4=30
GAMMAO=30
HO=<%= t - 2.0 %>
GAMMAU=0
HU=0
END_OF_PARAMS`,

		"bx": `ENDCUT_TYPE=37R
START_OF_PARAMS
NO_OF_PARAMS=13
E1=15
V1=90
V3=90
R1=0
A=20
L=80
D=0
GAMMA=0
HO=0
GAMMAU=0
HU=0
GAMMA2=0
HO2=0
END_OF_PARAMS`,
	}
)
