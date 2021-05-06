package mgf

var (
	AE = map[string]string{

		"DbS1": `ENDCUT_TYPE=3
START_OF_PARAMS
NO_OF_PARAMS=8
E1=0
A=15
V1=30
V3=0
E2=0
D=15
GAMMA=0
H=0
END_OF_PARAMS`,

		"AaL70R50": `ENDCUT_TYPE=5
START_OF_PARAMS
NO_OF_PARAMS=10
A=90
B=20
E1=0
R1=50
V1=90
E2=0
C=50
V3=0
GAMMA=0
H=0
END_OF_PARAMS`,

		"ba": `ENDCUT_TYPE=1
START_OF_PARAMS
NO_OF_PARAMS=5
V1=<%= 90.0 + data["angle"] %>
E1=0
E2=0
GAMMA=0
H=0
END_OF_PARAMS`,

		"babz": `ENDCUT_TYPE=23
START_OF_PARAMS
NO_OF_PARAMS=
E1=0
R1=0
V1=<%= 90.0 + data["angle"] %>
GAMMAF=<%= 90.0 - data["angle2"] %>
H=<%= data["h"] - t %>
GAMMAS=0
END_OF_PARAMS


ENDCUT_TYPE=9
START_OF_PARAMS
NO_OF_PARAMS=
E1=0
R1=0
V1=<%= 90.0 + data["angle"] %>
GAMMAF=<%= 90.0 - data["angle2"] %>
H=<%= data["h"] - t %>
GAMMAS=0
END_OF_PARAMS`,

		"bk": `ENDCUT_TYPE=4
START_OF_PARAMS
NO_OF_PARAMS=7
A=<%= h - data["h"] %>
E1=0
V1=<%= 90.0 + data["angle"] %>
V3=45
E2=0
H=0
GAMMA=0
END_OF_PARAMS`,

		"bs": `ENDCUT_TYPE=28
START_OF_PARAMS
NO_OF_PARAMS=12
E1=0
A=<%= data["l"] %>
B=<%= h - data["h"] %>
C=<%= data["snipe"] %>
V1=<%= 90.0 + data["angle"] %>
V3=90
R2=0.5
GAMMA=0
HO=0
GAMMAU=0
HU=0
D=10
END_OF_PARAMS`,
	}
)
