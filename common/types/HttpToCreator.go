package types

//HttpToCR This is the structure used by the communicate to CREATOR from HTTP module
type HttpToCR struct {
	Name    string //Name of the Instance
	NSlaves int    //Number of Slaves or Peeers
	Spec    WLSpec //IF there is any override of the basic worklaod spec (Additional json supplied while create)
}

//NewHttpToCR will create a structure for us
func NewHttpToCR(name string, nslaves int, payload string) HttpToCR {

	var H HttpToCR
	H.Name = name
	H.NSlaves = nslaves

	if payload != "" {
		H.Spec.FromJson(payload)
	} else {
		H.Spec.Copy(Cfg.WorkLoad)
	}

	return H
}
