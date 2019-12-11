// Copyright (C) 2019 <x6a@7n.io>
//
// pkg is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// pkg is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with pkg. If not, see <http://www.gnu.org/licenses/>.

package colors

import (
	"github.com/mgutz/ansi"
)

var Black = ansi.ColorFunc("black+bh")

var White = ansi.ColorFunc("white+bh")
var DarkWhite = ansi.ColorFunc("white+b")

var Blue = ansi.ColorFunc("blue+bh")
var DarkBlue = ansi.ColorFunc("blue+b")

var Cyan = ansi.ColorFunc("cyan+bh")
var DarkCyan = ansi.ColorFunc("cyan+b")

var Red = ansi.ColorFunc("red+bh")
var DarkRed = ansi.ColorFunc("red+b")

var Green = ansi.ColorFunc("gree+bh")
var DarkGreen = ansi.ColorFunc("gree+b")

var Magenta = ansi.ColorFunc("magenta+bh")
var DarkMagenta = ansi.ColorFunc("magenta+b")

var Yellow = ansi.ColorFunc("yellow+bh")
var DarkYellow = ansi.ColorFunc("yellow+b")
