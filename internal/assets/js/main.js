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
            const options = {
            chart: {
                type: 'area',
                height: 360,
                fontFamily: 'Inter, sans-serif',
                animations: {
                enabled: true,
                easing: 'easeinout',
                speed: 700
                }
            },

            series: [
                { name: 'Periodo Anterior', type: 'line', data: [] },
                { name: 'Periodo Actual', type: 'area', data: [] }
            ],

            colors: ['#94a3b8', '#4f46e5'],

            stroke: {
                curve: 'smooth',
                width: [2, 2],
                dashArray: [6, 0]
            },

            fill: {
                type: 'solid',
                opacity: [1, 0.25]   // línea sólida, área transparente elegante
            },

            markers: {
                size: [4, 4],       // línea + área con puntos
                strokeWidth: 2,
                hover: { size: 6 }
            },

            xaxis: {
                categories: [],
                tickPlacement: 'between',
                labels: {
                rotate: -45,
                style: {
                    fontSize: '11px',
                    colors: '#6b7280'
                }
                },
                axisBorder: { show: false },
                axisTicks: { show: false }
            },

            yaxis: {
                labels: {
                style: { colors: '#6b7280' }
                },
                title: {
                text: 'Cantidad de Atenciones',
                style: {
                    fontSize: '12px',
                    fontWeight: 500
                }
                }
            },

            grid: {
                borderColor: '#e5e7eb',
                strokeDashArray: 3,
                padding: {
                left: 12,
                right: 12
                }
            },

            legend: {
                position: 'top',
                horizontalAlign: 'center',
                markers: {
                width: 10,
                height: 10,
                radius: 12
                }
            },

            tooltip: {
                shared: true,
                intersect: false,
                theme: 'light'
            },

            dataLabels: { enabled: false },

            noData: {
                text: 'Seleccione un diagnóstico y rango de fechas',
                align: 'center',
                verticalAlign: 'middle',
                style: {
                color: '#6b7280',
                fontSize: '14px'
                }
            }
            };

            this.chart = new ApexCharts(this.$refs.chart, options);
            this.chart.render();
        },
        fetchData(id, inicio, fin) {
            console.log("Fetching data with:", { id, inicio, fin });
            if (!id || !inicio || !fin) {
                this.hasData = false;
                this.chart.updateOptions({ 
                    series: [
                        { name: 'Periodo Anterior', data: [] },
                        { name: 'Periodo Actual', data: [] }
                    ], 
                    xaxis: { categories: [] } 
                });
                return;
            }
            const params = new URLSearchParams({ IdDiagnostico: id, FechaInicio: inicio, FechaFin: fin });
            fetch(`/api/atenciones-por-dia?${params}`)
                .then(res => {
                    if (!res.ok) {
                        throw new Error(`HTTP error! status: ${res.status}`);
                    }
                    return res.json();
                })
                .then(data => {
                    console.log("Data received from API:", data);

                    if (data && data.periodoActual && data.periodoActual.length > 0) {
                        const fechasFormateadas = data.periodoActual.map(d => d.fecha);
                        const cantidadesActual = data.periodoActual.map(d => d.cantidadAtenciones);
                        const cantidadesAnterior = data.periodoAnterior.map(d => d.cantidadAtenciones);

                        console.log("Processed data for chart:", { fechasFormateadas, cantidadesActual, cantidadesAnterior });

                        this.chart.updateOptions({
                            series: [
                                { name: 'Periodo Anterior', data: cantidadesAnterior },
                                { name: 'Periodo Actual', data: cantidadesActual }
                            ],
                            xaxis: { categories: fechasFormateadas },
                            title: {
                                text: `Tendencia de Atenciones (${data.tipo})`,
                                align: 'left',
                                style: { fontSize: '16px', color: '#333' }
                            }
                        });
                        this.hasData = true;
                    } else {
                        console.log("No data received or data is empty. Resetting chart.");
                        this.hasData = false;
                        this.chart.updateOptions({ 
                            series: [
                                { name: 'Periodo Anterior', data: [] },
                                { name: 'Periodo Actual', data: [] }
                            ], 
                            xaxis: { categories: [] },
                            title: { text: '' }
                        });
                    }
                })
                .catch(error => {
                    console.error("Error fetching or processing data:", error);
                    this.hasData = false;
                    this.chart.updateOptions({ 
                        series: [
                            { name: 'Periodo Anterior', data: [] },
                            { name: 'Periodo Actual', data: [] }
                        ], 
                        xaxis: { categories: [] },
                        title: { text: 'Error al cargar datos' }
                    });
                });
        }
    }));

});
