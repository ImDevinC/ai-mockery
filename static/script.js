window.handleSubmit = async function(event) {
  event.preventDefault(); // Prevent the default form submission

  document.getElementById('result').innerHTML = '<p>Loading...</p>'; // Show loading message

  const form = document.getElementById('mockery-form');
  const formData = new FormData(form);
  const jsonData = {};
  formData.forEach((value, key) => {
    jsonData[key] = value;
  });

  try {
    const response = await fetch('/mockery', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(jsonData),
    });

    const result = await response.json();
    document.getElementById('result').innerHTML = `<pre>${JSON.stringify(result.insult, null, 2)}</pre>`;
  } catch (error) {
    document.getElementById('result').innerHTML = '<p>Error occurred while fetching data.</p>';
  }
}
