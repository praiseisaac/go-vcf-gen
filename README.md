Playing around with golang

# Go VCF File Generator
Generate vcf files from csv of names and phone numbers

### Running
`go run . path/to/csv [lines to skip]`

### Test Run
```bash
touch example.csv
echo 'first name, number' >> example.csv 
echo 'john doe, (123)-456-7890' >> example.csv
echo 'jane doe, (123)-456-7890' >> example.csv
go run . example.csv 1
```

### Improvements that could be made
1. Other file formats
2. Comma-separated names (first, last, number)
3. More interactivity
