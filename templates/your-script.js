fetch('/getdata')
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        console.log(response);
        return response.text();
    })
    .then(data => {
        // Log the data you received for debugging
        console.log(data);
        // Now try parsing it as JSON
     
     
    })
    .catch(error => {
        console.error('Error fetching data:', error);
    });