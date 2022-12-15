# Run training
# Each iteration runs a game and updates the weights

for i in {1..100}
do
    go test ./pkg/infra -run ^TestTrainQ$ -v infra
    python3 ../linreg.py
done