// internal/assets/js/main.js
document.addEventListener("alpine:init", () => {
  Alpine.data("sexoChartComponent", () => ({
    chart: null,
    init() {
      const options = {
        chart: {
          fontFamily: "Inter, sans-serif",
          type: "donut",
          height: 350,
          toolbar: { show: true },
          events: {
            dataPointMouseEnter: (event, chartContext, config) => {
              const seriesName = config.w.config.labels[config.dataPointIndex];
              const seriesValue = config.w.config.series[config.dataPointIndex];
              this.chart.updateOptions({
                plotOptions: { pie: { donut: { labels: { total: {
                  show: true,
                  label: seriesName,
                  formatter: () => seriesValue
                } } } } }
              });
            },
            dataPointMouseLeave: () => {
              this.chart.updateOptions({
                plotOptions: { pie: { donut: { labels: { total: {
                  show: true,
                  label: "Total",
                  formatter: (w) => w.globals.seriesTotals.reduce((a, b) => a + b, 0)
                } } } } }
              });
            }
          }
        },
        fill: {
          colors: ["#FF4560", "#2b7fff"]
        },
        series: [],
        labels: [],
        plotOptions: {
          pie: {
            donut: {
              size: "65%",
              labels: {
                show: true,
                total: {
                  show: true,
                  label: "Total",
                  formatter: (w) => {
                    if (w.globals.seriesTotals.length === 0)
                      return 0;
                    return w.globals.seriesTotals.reduce((a, b) => a + b, 0);
                  }
                }
              }
            }
          }
        },
        dataLabels: {
          enabled: true,
          formatter: function(val) {
            return val.toFixed(1) + "%";
          },
          dropShadow: { enabled: false }
        },
        legend: { position: "bottom" },
        noData: { text: "Seleccione un diagnóstico y rango de fechas." }
      };
      this.chart = new ApexCharts(this.$refs.chart, options);
      this.chart.render();
    },
    fetchData(id, inicio, fin) {
      if (!id || !inicio || !fin)
        return;
      const params = new URLSearchParams({ IdDiagnostico: id, FechaInicio: inicio, FechaFin: fin });
      fetch(`/api/sexo-por-diagnostico?${params}`).then((res) => res.json()).then((data) => {
        this.chart.updateOptions({
          series: data.map((d) => d.cantidadAtenciones),
          labels: data.map((d) => d.sexo)
        });
      });
    }
  }));
  Alpine.data("edadesChartComponent", () => ({
    chart: null,
    init() {
      const options = {
        chart: { type: "bar", height: 350, toolbar: { show: true }, fontFamily: "Inter, sans-serif" },
        series: [{ name: "Atenciones", data: [] }],
        yaxis: {
          categories: []
        },
        xaxis: {
          show: false,
          labels: {
            show: false
          }
        },
        plotOptions: {
          bar: {
            horizontal: false,
            dataLabels: { position: "top" },
            borderRadiusApplication: "end",
            borderRadius: 10
          }
        },
        fill: {
          colors: ["#00E396"]
        },
        dataLabels: {
          enabled: true,
          offsetY: -25,
          style: {
            fontSize: "14px",
            colors: ["#304758"]
          }
        },
        grid: {
          show: true,
          borderColor: "#e0e0e0",
          strokeDashArray: 2
        },
        noData: { text: "Seleccione un diagnóstico y rango de fechas." }
      };
      this.chart = new ApexCharts(this.$refs.chart, options);
      this.chart.render();
    },
    fetchData(id, inicio, fin) {
      if (!id || !inicio || !fin)
        return;
      const params = new URLSearchParams({ IdDiagnostico: id, FechaInicio: inicio, FechaFin: fin });
      fetch(`/api/edades-por-diagnostico?${params}`).then((res) => res.json()).then((data) => {
        this.chart.updateOptions({
          series: [{ data: data.map((d) => d.cantidadAtenciones) }],
          xaxis: { categories: data.map((d) => d.grupoEdad) }
        });
      });
    }
  }));
});
