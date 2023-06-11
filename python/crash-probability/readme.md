# 1

| odds | t   | y(0) == 0.05 |
|------|-----|--------------|
| 100  | 0.0 | 0.00597      |
| 101  | 0.1 | 0.00497      |
| 101  | 0.2 | 0.00897      |

```py
import pandas as pd

# Create a DataFrame
df = pd.DataFrame({
    'odds': odds,
    't': t,
    'y(0) == 0.05': y,  # Replace y with the actual list
})

# Write the DataFrame to an Excel file
df.to_excel('output.xlsx', index=False)
```