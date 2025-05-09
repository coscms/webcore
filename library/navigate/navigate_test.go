/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package navigate

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

func TestConfigDefaultsAsStore(t *testing.T) {
	assert.Equal(t, `"23434\"343434"`, strconv.Quote(`23434"343434`))
	v := List{
		{
			Name: `0`,
		},
		{
			Name: `1`,
		},
		{
			Name: `2`,
		},
	}
	v2 := List{
		{
			Name: `00`,
		},
		{
			Name: `01`,
		},
		{
			Name: `02`,
		},
	}

	v.Add(0, v2...)
	assert.Equal(t, List{
		{
			Name: `00`,
		},
		{
			Name: `01`,
		},
		{
			Name: `02`,
		},
		{
			Name: `0`,
		},
		{
			Name: `1`,
		},
		{
			Name: `2`,
		},
	}, v)

	v.Remove(0)
	assert.Equal(t, List{
		{
			Name: `01`,
		},
		{
			Name: `02`,
		},
		{
			Name: `0`,
		},
		{
			Name: `1`,
		},
		{
			Name: `2`,
		},
	}, v)

	v.Remove(4)
	assert.Equal(t, List{
		{
			Name: `01`,
		},
		{
			Name: `02`,
		},
		{
			Name: `0`,
		},
		{
			Name: `1`,
		},
	}, v)

	v.Set(0, List{
		{
			Name: `11`,
		},
		{
			Name: `12`,
		},
	}...)
	assert.Equal(t, List{
		{
			Name: `11`,
		},
		{
			Name: `12`,
		},
		{
			Name: `0`,
		},
		{
			Name: `1`,
		},
	}, v)
	v3 := List{
		{
			Name:   `11`,
			Action: `11`,
		},
		{
			Name:   `12`,
			Action: `12`,
		},
		{
			Name:   `0`,
			Action: `0`,
			Children: &List{
				{
					Name:   `1`,
					Action: `1`,
					Children: &List{
						{
							Name:   `2`,
							Action: `2`,
						},
					},
				},
			},
		},
	}
	assert.Equal(t, &Item{
		Name:   `2`,
		Action: `2`,
	}, v3.ChildItem(`0`, `1`, `2`))
	assert.Equal(t, &List{
		{
			Name:   `2`,
			Action: `2`,
		},
	}, v3.ChildList(`0`, `1`))
	v3.ChildList(`0`).ReplaceChild(`1`, `2`, &Item{
		Name:   `3`,
		Action: `3`,
	})
	assert.Equal(t, `[
  {
    "Name": "11",
    "Action": "11"
  },
  {
    "Name": "12",
    "Action": "12"
  },
  {
    "Name": "0",
    "Action": "0",
    "Children": [
      {
        "Name": "1",
        "Action": "1",
        "Children": [
          {
            "Name": "3",
            "Action": "3"
          }
        ]
      }
    ]
  }
]`, com.Dump(v3, true))
}

func TestGroup(t *testing.T) {
	g := Group{
		Group: `testG`,
		Label: `测试`,
	}
	RegisterGroup(`test`, g)
	assert.Equal(t, []Group{g}, navGroups[`test`].g)

	g2 := Group{
		Group: `testG2`,
		Label: `测试`,
	}
	RegisterGroup(`test`, g2)
	assert.Equal(t, []Group{g, g2}, navGroups[`test`].g)

	UnregisterGroup(`test`, `testG`)
	assert.Equal(t, []Group{g2}, navGroups[`test`].g)

	pn := NewProjectNavigates(`testCase`, `root`)
	pn.Add(Top, &List{
		{
			Display: true,
			Action:  `user`,
			Group:   ``,
			Children: &List{
				{
					Display: true,
					Action:  `index`,
					Group:   `user`,
				},
				{
					Display: true,
					Action:  `edit`,
					Group:   `user`,
				},
				{
					Display: true,
					Action:  `hidden`,
					Group:   ``,
				},
			},
		},
	})
	pn.Init()
	//com.Dump(navGroups)
	assert.Equal(t, []Group{
		Group{
			Group: `user`,
			Label: `User`,
		}, Group{
			Group: ``,
			Label: `Default`,
		},
	}, navGroups[`testCase.top.user`].g)
}
