document.addEventListener('alpine:init', () => {

    Alpine.data('sexoChartComponent', () => ({
        chart: null,
        hasData: false,
        init() {
            const options = {
                chart: {
                    fontFamily: 'Inter, sans-serif',
                    type: 'donut', 
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
                                }}}}}
                            });
                        },
                        dataPointMouseLeave: () => {
                            this.chart.updateOptions({
                                plotOptions: { pie: { donut: { labels: { total: {
                                    show: true,
                                    label: 'Total',
                                    formatter: (w) => w.globals.seriesTotals.reduce((a, b) => a + b, 0)
                                }}}}}
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
                            size: '65%',
                            labels: {
                                show: true,
                                total: {
                                    show: true,
                                    label: 'Total',
                                    formatter: (w) => {
                                        if (w.globals.seriesTotals.length === 0) return 0;
                                        return w.globals.seriesTotals.reduce((a, b) => a + b, 0)
                                    }
                                }
                            }
                        }
                    }
                },
                dataLabels: {
                    enabled: true,
                    formatter: function (val) {
                        return val.toFixed(1) + '%'
                    },
                    dropShadow: { enabled: false },
                },
                legend: { position: 'bottom' },
                noData: { text: 'Seleccione un diagnóstico y rango de fechas.' }
            };
            this.chart = new ApexCharts(this.$refs.chart, options);
            this.chart.render();
        },
        fetchData(id, inicio, fin) {
            if (!id || !inicio || !fin) {
                this.hasData = false;
                this.chart.updateOptions({ series: [], labels: [] });
                return;
            }
            const params = new URLSearchParams({ IdDiagnostico: id, FechaInicio: inicio, FechaFin: fin });
            fetch(`/api/sexo-por-diagnostico?${params}`)
                .then(res => res.json())
                .then(data => {
                    if(data && data.length > 0) {
                        this.chart.updateOptions({
                            series: data.map(d => d.cantidadAtenciones),
                            labels: data.map(d => d.sexo)
                        });
                        this.hasData = true;
                    } else {
                        this.hasData = false;
                        this.chart.updateOptions({ series: [], labels: [] });
                    }
                });
        }
    }));

    // Jeferson Rivas - 205
    Alpine.data('edadesChartComponent', () => ({
        chart: null,
        hasData: false,
        init() {
            const options = {
                chart: { type: 'bar', height: 350, toolbar: { show: true }, fontFamily: 'Inter, sans-serif', },
                series: [{ name: 'Atenciones', data: [] }],
                xaxis: {
                    categories: []
                },
                yaxis: {
                    show: false,
                    labels: {
                        show: false
                    }
                },
                plotOptions: {
                    bar: {
                        horizontal: false,
                        dataLabels: { position: 'top' },
                        borderRadiusApplication: 'end', // Apply to end of the bars
                        borderRadius: 10, // Rounded borders
                    }
                },
                fill: {
                    colors: ["#00E396"]
                },
                dataLabels: {
                    enabled: true,
                    offsetY: -25,
                    style: {
                        fontSize: '14px',
                        colors: ["#304758"]
                    }
                },
                grid: {
                    show: true,
                    borderColor: '#e0e0e0',
                    strokeDashArray: 2,
                },
                noData: { text: 'Seleccione un diagnóstico y rango de fechas.' }
            };
            this.chart = new ApexCharts(this.$refs.chart, options);
            this.chart.render();
        },
        fetchData(id, inicio, fin) {
            if (!id || !inicio || !fin) {
                this.hasData = false;
                this.chart.updateOptions({ series: [{ data: [] }], xaxis: { categories: [] } });
                return;
            }
            const params = new URLSearchParams({ IdDiagnostico: id, FechaInicio: inicio, FechaFin: fin });
            fetch(`/api/edades-por-diagnostico?${params}`)
                .then(res => res.json())
                .then(data => {
                    if (data && data.length > 0) {
                        this.chart.updateOptions({
                            series: [{ data: data.map(d => d.cantidadAtenciones) }],
                            xaxis: { categories: data.map(d => d.grupoEdad) }
                        });
                        this.hasData = true;
                    } else {
                        this.hasData = false;
                        this.chart.updateOptions({ series: [{ data: [] }], xaxis: { categories: [] } });
                    }
                });
        }
    }));

    Alpine.data('distritosChartComponent', () => ({
        chart: null,
        hasData: false,
        init() {
            const options = {
                chart: { type: 'bar', height: 350, toolbar: { show: true }, fontFamily: 'Inter, sans-serif', },
                series: [{ name: 'Atenciones', data: [] }],
                xaxis: {
                    categories: [],
                    labels: {
                        rotate: -45,
                        style: {
                            fontSize: '11px'
                        }
                    }
                },
                plotOptions: {
                    bar: {
                        horizontal: false,
                        dataLabels: { position: 'top' },
                        borderRadiusApplication: 'end',
                        borderRadius: 10,
                    }
                },
                fill: {
                    colors: ["#2563eb"]
                },
                dataLabels: {
                    enabled: true,
                    offsetY: -25,
                    style: {
                        fontSize: '14px',
                        colors: ["#304758"]
                    }
                },
                grid: {
                    show: true,
                    borderColor: '#e0e0e0',
                    strokeDashArray: 2,
                },
                noData: { text: 'Seleccione un diagnóstico y rango de fechas.' }
            };
            this.chart = new ApexCharts(this.$refs.chart, options);
            this.chart.render();
        },
        fetchData(id, inicio, fin) {
            if (!id || !inicio || !fin) {
                this.hasData = false;
                this.chart.updateOptions({ series: [{ data: [] }], xaxis: { categories: [] } });
                return;
            }
            const params = new URLSearchParams({ IdDiagnostico: id, FechaInicio: inicio, FechaFin: fin });
            fetch(`/api/distritos-por-diagnostico?${params}`)
                .then(res => res.json())
                .then(data => {
                    if (data && data.length > 0) {
                        const top10 = data.slice(0, 10);
                        this.chart.updateOptions({
                            series: [{ data: top10.map(d => d.cantidadAtenciones) }],
                            xaxis: { categories: top10.map(d => `${d.nombreDistrito}`) }
                        });
                        this.hasData = true;
                    } else {
                        this.hasData = false;
                        this.chart.updateOptions({ series: [{ data: [] }], xaxis: { categories: [] } });
                    }
                });
        }
    }));

    Alpine.data('tiempoChartComponent', () => ({
        chart: null,
        hasData: false,
        init() {
            const mesesCortos = ['Ene', 'Feb', 'Mar', 'Abr', 'May', 'Jun', 'Jul', 'Ago', 'Sep', 'Oct', 'Nov', 'Dic'];

            const options = {
                chart: {
                    type: 'area',
                    height: 350,
                    toolbar: { show: true },
                    fontFamily: 'Inter, sans-serif',
                    zoom: { enabled: true }
                },
                series: [{ name: 'Atenciones', data: [] }],
                xaxis: {
                    categories: [],
                    labels: {
                        rotate: -45,
                        style: {
                            fontSize: '11px'
                        }
                    }
                },
                yaxis: {
                    title: {
                        text: 'Cantidad de Atenciones'
                    }
                },
                dataLabels: {
                    enabled: false
                },
                stroke: {
                    curve: 'smooth',
                    width: 2
                },
                fill: {
                    colors: ['#2563eb']
                },
                markers: {
                    size: 4,
                    colors: ['#2563eb'],
                    strokeColors: '#fff',
                    strokeWidth: 2,
                    hover: {
                        size: 6
                    }
                },
                grid: {
                    show: true,
                    borderColor: '#e0e0e0',
                    strokeDashArray: 2,
                },
                tooltip: {
                    custom: function({ series, seriesIndex, dataPointIndex, w }) {
                        const value = series[seriesIndex][dataPointIndex];
                        const fecha = w.globals.labels[dataPointIndex];
                        return '<div class="px-3 py-2 bg-white border border-slate-200 rounded shadow-lg">' +
                               '<div class="text-xs text-slate-500">' + fecha + '</div>' +
                               '<div class="text-sm font-semibold text-slate-800">Atenciones: ' + value + '</div>' +
                               '</div>';
                    }
                },
                noData: { text: 'Seleccione un diagnóstico y rango de fechas.' }
            };
            this.chart = new ApexCharts(this.$refs.chart, options);
            this.chart.render();
            this.mesesCortos = mesesCortos;
        },
        fetchData(id, inicio, fin) {
            if (!id || !inicio || !fin) {
                this.hasData = false;
                this.chart.updateOptions({ series: [{ data: [] }], xaxis: { categories: [] } });
                return;
            }
            const params = new URLSearchParams({ IdDiagnostico: id, FechaInicio: inicio, FechaFin: fin });
            fetch(`/api/atenciones-por-dia?${params}`)
                .then(res => res.json())
                .then(data => {
                    if (data && data.length > 0) {
                        const fechasFormateadas = data.map(d => {
                            const [year, month, day] = d.fecha.split('T')[0].split('-');
                            return `${day} ${this.mesesCortos[parseInt(month) - 1]}`;
                        });

                        const cantidades = data.map(d => d.cantidadAtenciones);

                        this.chart.updateOptions({
                            series: [{ name: 'Atenciones', data: cantidades }],
                            xaxis: { categories: fechasFormateadas }
                        });
                        this.hasData = true;
                    } else {
                        this.hasData = false;
                        this.chart.updateOptions({ series: [{ data: [] }], xaxis: { categories: [] } });
                    }
                });
        }
    }));
});
