package main


type FuelInput struct {
	Hydrogen float64
	Carbon   float64
	Sulfur   float64
	Nitrogen float64
	Oxygen   float64
	Moisture float64
	Ash      float64
}


type FuelOilInput struct {
	Carbon         float64
	Hydrogen       float64
	Sulfur         float64
	Vanadium       float64
	Oxygen         float64
	Moisture       float64
	Ash            float64
	HeatCombustion float64
}

type DryMassResult struct {
	K float64
	H float64
	C float64
	S float64
	N float64
	O float64
	A float64
}


type CombustibleMassResult struct {
	K float64
	H float64
	C float64
	S float64
	N float64
	O float64
}

type HeatCombustionResult struct {
	Q            float64
	QDry         float64
	QCombustible float64
}

type FuelOilCompositionResult struct {
	H float64
	C float64
	S float64
	V float64
	A float64
	O float64
}


func calculateDryMass(input *FuelInput) *DryMassResult {
	k := 100.0 / (100.0 - input.Moisture)
	return &DryMassResult{
		K: k,
		H: k * input.Hydrogen,
		C: k * input.Carbon,
		S: k * input.Sulfur,
		N: k * input.Nitrogen,
		O: k * input.Oxygen,
		A: k * input.Ash,
	}
}


func calculateCombustibleMass(input *FuelInput) *CombustibleMassResult {
	k := 100.0 / (100.0 - input.Moisture - input.Ash)
	return &CombustibleMassResult{
		K: k,
		H: k * input.Hydrogen,
		C: k * input.Carbon,
		S: k * input.Sulfur,
		N: k * input.Nitrogen,
		O: k * input.Oxygen,
	}
}

func calculateFuelHeatCombustion(input *FuelInput) *HeatCombustionResult {
	q := (339*input.Carbon + 1030*input.Hydrogen - 108.8*(input.Oxygen-input.Sulfur) - 25*input.Moisture) / 1000
	qDry := (q + 0.025*input.Moisture) * (100 / (100 - input.Moisture))
	qCombustible := (q + 0.025*input.Moisture) * (100 / (100 - input.Moisture - input.Ash))
	return &HeatCombustionResult{
		Q:            q,
		QDry:         qDry,
		QCombustible: qCombustible,
	}
}

func calculateFuelOilComposition(input *FuelOilInput) *FuelOilCompositionResult {
	factor := (100 - input.Moisture - input.Ash) / 100
	factorW := (100 - input.Moisture) / 100
	return &FuelOilCompositionResult{
		H: input.Hydrogen * factor,
		C: input.Carbon * factor,
		S: input.Sulfur * factor,
		V: input.Vanadium * factorW,
		A: input.Ash * factorW,
		O: input.Oxygen * factor,
	}
}


func calculateFuelOilHeatCombustion(input *FuelOilInput) float64 {
	return input.HeatCombustion*((100-input.Moisture-input.Ash)/100) - 0.025*input.Moisture
}
