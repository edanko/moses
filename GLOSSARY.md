# Glossary

Terms used by the Moses cutting pipeline. Definitions describe the business
meaning, not the implementation.

## Project

A shipbuilding project or order number that owns sections, parts and remnants.
Example: `056001`.

## Section

An assembly section or block within a project. A section groups parts that are
cut and reported together. Also called a block.

## Profile designation

The normalized name of a steel profile, for example `RP160X10`, `FB80X8` or
`AE250X90X10X15`. It identifies the cross-section family and its dimensions.
Not the same as a part: many parts can use the same designation.

## Profile family

The letter prefix of a profile designation: `RP` (bulb flat), `FB` (flat bar)
or `AE` (angle). Configuration may define defaults per family.

## Steel grade

The material grade of the profile, for example `A40` or `S235`. The old code
called this "quality"; the domain now calls it grade consistently.

## Part

One production piece to be cut from a stock bar. A part has a project, section,
position number, profile designation, steel grade, length, quantity, two end
cuts and any internal cutouts. "Profile" in older code meant this part; it now
means the profile designation only.

## Position number

The drawing position or mark of a part within its section. Used to match a cut
part back to the drawing. The old code abbreviated this as `posNo`.

## End cut

The machining applied to the left or right end of a part. An end cut has a kind
(for example `21`) and numeric parameters such as `v1`, `v2` or `snipe`.

## Cutout

An internal cut in a part: a hole, notch or other `icut`. A cutout has a kind
and numeric parameters such as `x`, `y` and `r`.

## Stock bar

The material a nest is cut from. A stock bar has a profile designation, steel
grade and length. It can be a new bar, a reusable remnant, or an explicit
over-length bar created when no configured stock fits.

## Remnant

An unused offcut left after cutting a nest. A remnant can be reused as a stock
bar for a later nest. Remnants are stored in the remnant CSV with an integer id,
source, marking and consumption state.

## Nest

One stock bar together with the parts cut from it. A nest has one stock bar,
one or more parts and optionally the remnant it consumes.

## Effective length

The length a part occupies in a nest: its raw length plus the spacing required
by its left and right end cuts.

## Spacing

The allowance added at an end cut so the part can be machined. Spacing is
resolved per profile and end-cut kind, with family and global defaults.

## Material utilization

The percentage of stock length used by parts. The old reports called this
"nesting percent".
