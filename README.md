# fragility-curves
A fragility curve plugin for cloud wat.

The fragility curve plugin takes a set of input seeds representing natural variability and knowledge uncertainty for a given event, and will sample all fragility curves with a curve sample using the knowledge uncertainty seed, and a value sample (failure elevation) using the event seed.

## Example Fragility Curve
The table below gives an example of a fragility curve, the x values are probability and the y values are elevation. The elevation is represented as a distribution for a given probability of failure of the infrastructure being modeled. The elevation is expressed as a distributed variable to represent the uncertainty in the predicted elevation of failure for any given probability.

[Image of Example Fragility Curve](example_fragility_curve.png)

|Probability|Min Elevation|Most Likely Elevation|Max Elevation|Distribution Type|
|---|---|---|---|---|
|0.0|98|99|100|Triangular|
|0.1|99|100|101|Triangular|
|0.2|100|101|102|Triangular|
|0.3|101|102|103|Triangular|
|0.4|102|103|104|Triangular|
|0.5|103|104|105|Triangular|
|0.6|104|105|106|Triangular|
|0.7|105|106|107|Triangular|
|0.8|106|107|108|Triangular|
|0.9|107|108|109|Triangular|
|1.0|108|109|110|Triangular|

That same table is represented in this example json data file
```json
{
    "name":"testModel",
	"locations": [{
		"location": "levee1",
		"probability-stage": {
			"xvalues": [0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1],
			"ydistributions": [{
				"type": "TriangularDistribution",
				"parameters": {
					"min": 98,
					"mostlikely": 99,
					"max": 100
				}
			}, {
				"type": "TriangularDistribution",
				"parameters": {
					"min": 99,
					"mostlikely": 100,
					"max": 101
				}
			}, {
				"type": "TriangularDistribution",
				"parameters": {
					"min": 100,
					"mostlikely": 101,
					"max": 102
				}
			}, {
				"type": "TriangularDistribution",
				"parameters": {
					"min": 101,
					"mostlikely": 102,
					"max": 103
				}
			}, {
				"type": "TriangularDistribution",
				"parameters": {
					"min": 102,
					"mostlikely": 103,
					"max": 104
				}
			}, {
				"type": "TriangularDistribution",
				"parameters": {
					"min": 103,
					"mostlikely": 104,
					"max": 105
				}
			}, {
				"type": "TriangularDistribution",
				"parameters": {
					"min": 104,
					"mostlikely": 105,
					"max": 106
				}
			}, {
				"type": "TriangularDistribution",
				"parameters": {
					"min": 105,
					"mostlikely": 106,
					"max": 107
				}
			}, {
				"type": "TriangularDistribution",
				"parameters": {
					"min": 106,
					"mostlikely": 107,
					"max": 108
				}
			}, {
				"type": "TriangularDistribution",
				"parameters": {
					"min": 107,
					"mostlikely": 108,
					"max": 109
				}
			}, {
				"type": "TriangularDistribution",
				"parameters": {
					"min": 108,
					"mostlikely": 109,
					"max": 110
				}
			}]
		}
	}]
}
```
