package mgf

var (
	FB = map[string]string{

		"bck2x45-8": `ENDCUT_TYPE=13
START_OF_PARAMS
NO_OF_PARAMS=10
E1=<%= data["b"] %>
B=<%= data["h"] %>
V1=<%= 90.0 + data["angle"] %>
V2=<%= 90.0 - data["angle2"] %>
V3=90
GAMMA=45
HO=<%= t / 2.0 - 4.0 %>
GAMMAU=45
HU=8
END_OF_PARAMS`,

		"bk2x45-8": `ENDCUT_TYPE=11
START_OF_PARAMS
NO_OF_PARAMS=6
V1=<%= 90.0 + data["angle"] %>
E1=<%= data["h"] %>
GAMMA=45
HO=<%= t / 2.0 - 4.0 %>
GAMMAU=45
HU=8
END_OF_PARAMS`,

		"bl": `ENDCUT_TYPE=11R
START_OF_PARAMS
NO_OF_PARAMS=6
V1=<%= 90.0 + data["angle"] %>
R1=<%= data["r"] %>
R2=0
GAMMA=0
HO=0
GAMMAU=0
HU=0
END_OF_PARAMS`,

		"bv": `ENDCUT_TYPE=EF06
START_OF_PARAMS
NO_OF_PARAMS=6
A=<%= a(data["angle"]) %>
B=<%= data["h"] %>
C=<%= data["b"] %>
R1=0
R2=0
R3=<%= data["r"] %>
END_OF_PARAMS`,

		"bvf": `ENDCUT_TYPE=EF06
START_OF_PARAMS
NO_OF_PARAMS=6
A=<%= a(data["angle"]) %>
B=<%= data["h"] %>
C=<%= data["b"] %>
R1=0
R2=0
R3=<%= data["r"] %>
END_OF_PARAMS`,

		"bvg": `ENDCUT_TYPE=EF06
START_OF_PARAMS
NO_OF_PARAMS=6
A=<%= a(data["angle"]) %>
B=<%= data["h"] %>
C=<%= data["b"] %>
R1=0
R2=0
R3=<%= data["r"] %>
END_OF_PARAMS`,

		"cx": `ENDCUT_TYPE=14
START_OF_PARAMS
NO_OF_PARAMS=8
B=<%= data["h"] %>
C=<%= (h - data["h"]) / 2.0 %>
V1=90
V2=<%= 90.0 - data["angle"] %>
GAMMA=0
HO=0
GAMMAU=0
HU=0
END_OF_PARAMS`,

		"cx_dok": `ENDCUT_TYPE=14
START_OF_PARAMS
NO_OF_PARAMS=8
B=<%= 3.0 * t %>
C=<%= (h - (3.0 * t)) / 2.0 %>
V1=90
V2=<%= cx_dok(data["l"]) %>
GAMMA=0
HO=0
GAMMAU=0
HU=0
END_OF_PARAMS`,
	}
)
