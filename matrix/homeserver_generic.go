package matrix

func ReadHomeserverFile(homeserverType HomeserverType, basePath string) {
	switch homeserverType {
	case SynapseType:
		readSynapseFile(basePath)
	case DendriteType:
		// TODO: Dendrite support
	}
}

func WriteHomeserverFile(homeserverType HomeserverType, basePath string) {
	switch homeserverType {
	case SynapseType:
		writeSynapseFile(basePath)
	case DendriteType:
		// TODO: Dendrite support
	}
}
