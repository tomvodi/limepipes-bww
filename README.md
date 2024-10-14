# limepipes-plugin-bww

A LimePipes plugin that handles Bagpipe Music Writer Gold and Bagpipe Player files.

## Project Structure

This project contains the implementation for parsing Bagpipe Music Writer Gold files and Bagpipe Player files 
and translate them into the music model defined in the [LimePipes plugin API](https://github.com/tomvodi/limepipes-plugin-api).

`.bww` and `.bmw` files don't have a well defined structure and the Bagpipe Player is very fault tolerant with the input files.
It is undefined, if a timeline end symbols should be in a staff or outside.
Logically, the symbols should be inside the staff, but the Bagpipe Player is able to handle both cases.

This plugin tokenizes `.bww` and `.bmw` files into tokens that make sense for this specific file format.
These tokens are then translated into a intermediate representation (filestructure) that is then translated 
into the music model. This makes it possible to fix issues in the input files before translating them into 
the music model.

As first step, the parser takes the input file and tokenizes it into a list of tokens that represent the file structure.
This doesn't include the specific musical symbols like notes, rests, etc. but the structure of the file with tune header 
fields, measures and the symbols in general. These general symbols contain only their corresponding symbol text and position.

In the following step, the parser translates the tokens into the music model. Here, the symbol text is translated into the
corresponding music model symbol with the help of the symbol mapper. Here also happens the merging of symbols that belong together.
As embellishments and melody notes are two symbols in the Bagpipe Player file format, they are merged into one symbol in the music model.
This is also true for the melody note dots and other.

### Fixing the input files

Bagpipe Player files don't have the ability to specify an arranger. Most of the time the arranger is specified in the composer field 
as "Charly Composer, arr. Willi World" or similar. This plugin tries to extract the arranger from the composer field.
The same is true for "Traditional" tunes, where the composer field contains something like "Traditional, ..." or "Trad.".

Some other fixes are:
- Remove time signature from the tune type field (6/8 March -> March)
- Capitalize the tune type (reel -> Reel)
- Trim spaces from fields
- remove underscores from the title field
- Fix cases for title (my tune -> My Tune)

### Directory Structure

`bww`

The directory for all the parser related stuff.

`bwwfile`

Here is all code related to the intermediate file structure parsing that is later used by the parser/converter itself.

`pluginimplementation`

The implementation of the LimePipes plugin that utilizes the parser and the file structure.

## Build

`go build ./...` builds the plugin.

