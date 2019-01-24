package catalog

import (
	"fmt"
	"testing"
)

func TestModifyParts(t *testing.T) {
	parameterIDToChange := "ac-10_prm_3"
	parameterIDNotToChange := "ac-10_prm_2"
	parameterVal := "[CHANGED]"
	shouldChange := fmt.Sprintf(`this should change. <insert param-id="%s">`, parameterIDToChange)
	shouldNotChange := fmt.Sprintf(`this should not change <insert param-id="%s">`, parameterIDNotToChange)
	afterChange := fmt.Sprintf(`this should change. %s`, parameterVal)
	prose := Prose{
		P: []P{
			P{
				Raw: shouldChange,
			},
		},
	}
	nestedProse := Prose{
		P: []P{
			P{
				Raw: shouldChange,
			},
		},
	}
	c := Catalog{
		Groups: []Group{
			Group{
				Controls: []Control{
					Control{
						Parts: []Part{
							Part{
								Prose: &prose,
								Parts: []Part{
									Part{
										Prose: &Prose{
											P: []P{
												P{
													Raw: shouldNotChange,
												},
											},
										},
										Parts: []Part{
											Part{
												Prose: &nestedProse,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	c.Groups[0].Controls[0].Parts[0].ModifyProse(parameterIDToChange, parameterVal)

	if c.Groups[0].Controls[0].Parts[0].Prose.P[0].Raw != afterChange {
		t.Error("part not modified")
	}

	if c.Groups[0].Controls[0].Parts[0].Parts[0].Prose.P[0].Raw != shouldNotChange {
		t.Error("part got modified which shouldnt")
	}
	if c.Groups[0].Controls[0].Parts[0].Parts[0].Parts[0].Prose.P[0].Raw != afterChange {
		t.Error("part not modified")
	}
}
