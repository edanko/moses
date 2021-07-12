# moses

Shipbuilding production profile part nesting tool. You can take as input following formats:

- Simple txt format
- AVEVA Generic file (*.gen) from profile interface
- AVEVA csv&lst files from "Profile sketch & list"
- Production drawings from Cadmatic (in dxf format)

and output nested bars in mgf file format (used in AutoCAM MOSES software)

# TODO
- [ ] read weight form *.gen and store in part model?
- [ ] rewrite this to restful backend (gofiber) and vue.js frontend
- [ ] store remnants in mongodb
- [ ] store profile in mongodb
- [ ] add second profile cutting line output
- [ ] add separate program to import parts to db
- [ ] validate everything
- [ ] some checker program (files vs db) or integrate this to importer