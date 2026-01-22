# MarkView Theme Showcase

Welcome to **MarkView**! This document demonstrates the beautiful *custom themes* with syntax highlighting.

## Color Scheme Features

MarkView comes with two gorgeous themes:

### Dark Theme (Dracula-inspired)
- Deep purple-gray background
- Vibrant syntax highlighting
- Easy on the eyes for extended reading

### Light Theme (GitHub-inspired)
- Clean white background
- Professional color palette
- Perfect for daytime reading

## Text Formatting

You can use **bold text**, *italic text*, and even ***bold italic text*** for emphasis.

Here's some `inline code` that looks great in both themes.

## Code Blocks with Syntax Highlighting

### Python Example

```python
def fibonacci(n):
    """Calculate the nth Fibonacci number."""
    if n <= 1:
        return n
    return fibonacci(n - 1) + fibonacci(n - 2)

# Test the function
result = fibonacci(10)
print(f"The 10th Fibonacci number is: {result}")
```

### Go Example

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Print current time
    now := time.Now()
    fmt.Printf("Current time: %s\n", now.Format(time.RFC3339))

    // Calculate factorial
    n := 5
    factorial := 1
    for i := 1; i <= n; i++ {
        factorial *= i
    }
    fmt.Printf("Factorial of %d is %d\n", n, factorial)
}
```

### JavaScript Example

```javascript
// Async function with arrow syntax
const fetchUserData = async (userId) => {
    try {
        const response = await fetch(`/api/users/${userId}`);
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching user:', error);
        throw error;
    }
};

// Class example
class Calculator {
    add(a, b) {
        return a + b;
    }

    multiply(a, b) {
        return a * b;
    }
}

const calc = new Calculator();
console.log(calc.add(5, 3)); // Output: 8
```

### Rust Example

```rust
// Struct with methods
struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    fn area(&self) -> u32 {
        self.width * self.height
    }

    fn can_hold(&self, other: &Rectangle) -> bool {
        self.width > other.width && self.height > other.height
    }
}

fn main() {
    let rect1 = Rectangle { width: 30, height: 50 };
    let rect2 = Rectangle { width: 10, height: 40 };

    println!("Area: {}", rect1.area());
    println!("Can hold: {}", rect1.can_hold(&rect2));
}
```

## Lists

### Unordered Lists

- First item with **bold text**
- Second item with *italic text*
- Third item with `inline code`
  - Nested item one
  - Nested item two
  - Nested item three

### Ordered Lists

1. Configure your environment
2. Install dependencies
3. Run the application
4. Enjoy the beautiful themes!

## Blockquotes

> "The best way to predict the future is to invent it."
>
> — Alan Kay

> This is a longer blockquote that spans multiple lines.
> It demonstrates how blockquotes are rendered with the custom styling.
> Notice the nice vertical bar and italic text!

## Links

Check out these resources:
- [Fyne GUI Toolkit](https://fyne.io)
- [Goldmark Parser](https://github.com/yuin/goldmark)
- [Chroma Syntax Highlighter](https://github.com/alecthomas/chroma)

## Theme Switching

To switch between themes:
1. Go to **View** menu
2. Select **Light Theme** or **Dark Theme**
3. The entire interface updates instantly!

## Tables (GFM)

| Feature | Dark Theme | Light Theme |
|---------|------------|-------------|
| Background | Purple-Gray | Pure White |
| Text Color | Off-White | Dark Gray |
| Syntax Highlighting | Vibrant | Professional |
| Code Background | Darker | Light Gray |
| Link Color | Cyan | Blue |

## Conclusion

MarkView's custom themes provide a beautiful, distraction-free environment for reading and writing markdown. The syntax highlighting makes code blocks pop, while the carefully chosen colors reduce eye strain.

**Try switching themes now to see the difference!**

---

*Made with ❤️ using Go and Fyne*
