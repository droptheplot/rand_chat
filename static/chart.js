window.onload = function(){
  fetch('/api/chart')
    .then(res => res.json())
    .then((out) => {
      new Chart(document.getElementById("chart"), {
        type: 'line',
        data: {
          labels: out["dates"],
          datasets: [{
            data: out["counts"],
            backgroundColor: 'rgba(0, 123, 255, 0.2)',
            borderColor: 'rgba(0, 123, 255, 1)',
            borderWidth: 1
          }],
        },
        options: {
          legend: {
            display: false,
          },
          tooltips: {
            enabled: false,
          },
          elements: {
            point: {
              radius: 0,
            }
          },
          scales: {
            yAxes: [{
              gridLines: {
                tickMarkLength: 0,
              },
              ticks: {
                fontColor: "#aaa",
                padding: 10,
              },
            }],
            xAxes: [{
              gridLines: {
                tickMarkLength: 0,
              },
              ticks: {
                fontColor: "#aaa",
                padding: 10,
              },
            }],
          },
          maintainAspectRatio: false,
        },
      });
    });
};
