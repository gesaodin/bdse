
var _TIPO = 0; //Tipo de Seleccion
var opciones = {
    "paging": true,
    "ordering": true,
    "info": true,
    "searching": true,
    "language": {
        "decimal": ",",
        "thousands": "."
    }
};
//$('#reporte').DataTable(opciones);
$('#reporteSaldos').DataTable(opciones);


$(function () {
    /**
     * Declaración de variables globales
     */
    //var t = $('#reporte').DataTable();

    $(".select2").select2();
    var rS = $('#reporteSaldos').DataTable();
    $("[data-mask]").inputmask();
    if ($('#cuentahaber').val() != undefined) LCuentaM();
    if ($('#estado').val() != undefined) LEstado();
    if ($('#taquilla').val() != undefined) LProgramas();
    if($('#mdlER').html() != undefined) LCuentaB();
    CargarPerfil();
    CargarCalendario();

    $("#cargando").hide();
    socket.addEventListener('message', function (e) {
        var js = JSON.parse(e.data);
        if (js.tipo != null) {
            switch (js.tipo) {
                case 0:
                    break;
                case 1:
                    CrearNotificacion(js.tiempo, js.msj)
                    break;
                case 2:
                    break;
                case 3:
                    ActivarChat(js.De, js.msj, js.tiempo);
                case 33:
                    //Error en los archivos enviados
                    CNErr(js.tiempo, js.msj);
                default:
                    break;
            }
            //console.log(e.data);
        }

    });
    socket.addEventListener('open', function (e) {
        console.log("Se establecio la conexión con el socket...");
    });
    socket.addEventListener('close', function (e) {
        console.log("Se ha cerrado la conexión");
    });
    if ($("#listagrupo").html() != undefined)LGrupos();

});


function CargarPerfil() {
    if (sessionStorage.perfil == null) {
        var data = JSON.stringify({ id: 1 });
        $.post("api/perfil/comercializadora", data)
            .done(function (data) {
                sessionStorage.setItem('perfil', JSON.stringify(data));
                $("#lblGastos").html(parseFloat(data.gastos).toFixed(2) + " Bs.");
            });
    } else {
        if ($("#lblGastos").html() != undefined) {
            perfil = JSON.parse(sessionStorage.perfil);
            $("#lblGastos").html(parseFloat(perfil.gastos).toFixed(2) + " Bs.");
            $.each(perfil.lgrupo, function(c, v){
                console.log(v.nombre);
            });
        }
    }
}

function LGrupos(){
    perfil = JSON.parse(sessionStorage.perfil);
    $.each(perfil.lgrupo, function(c, v){
        $("#listagrupo").append(ListarGrupos(v.nombre));
    });
}

/**
 * Listar Cuentas para Movimientos
 */
function LCuentaM(id) {
    $('#cuentadebe').html('<option value="--" >Seleccionar...</option>');
    $('#cuentahaber').html('<option value="--" >Seleccionar...</option>');
    var data = JSON.stringify({operacion:0});
    $.post("api/movimiento/listarcuentas", data)
        .done(function (data) {
            $.each(data, function (c, v) {
                $('#cuentadebe').append('<option value="' + v.oid + '">' + v.oid + ' | ' + v.nombre + '</option>');
                $('#cuentahaber').append('<option value="' + v.oid + '">' + v.oid + ' | ' + v.nombre + '</option>')
            })
        });
}
/**
 * Listar Cuentas para Bancos
 */
function LCuentaB() {

    $('#cuentaer').html('<option value="--" >Seleccionar...</option>');
    var data = JSON.stringify({operacion:1});
    $.post("api/movimiento/listarcuentas", data)
        .done(function (data) {
            $.each(data, function (c, v) {
                $('#cuentaer').append('<option value="' + v.oid + '">' + v.oid + ' | ' + v.nombre + '</option>');
            })
        });
}



/**
 * Listar Estados
 */
function LEstado() {
    var data = JSON.stringify({ id: 0 });
    $('#estado').html('<option value="--" >Seleccionar...</option>');
    $('#ciudad').html('<option value="--" >Seleccionar...</option>');
    $('#parroquia').html('<option value="--" >Seleccionar...</option>');
    $('#municipio').html('<option value="--" >Seleccionar...</option>');
    $.post("api/localizacion/consultarestado", data)
        .done(function (data) {
            $.each(data, function (c, v) {
                $('#estado').append('<option value="' + v.id + '">' + v.nombre + '</option>');
            })
        });
}

/**
 * Listar Ciudad
 */
function LCiudad() {
    var id = $('#estado option:selected').val();
    var data = JSON.stringify({ ide: parseInt(id) });
    $('#ciudad').html('<option value="--" >Seleccionar...</option>');
    $.post("api/localizacion/consultarciudad", data)
        .done(function (data) {
            $.each(data, function (c, v) {
                $('#ciudad').append('<option value="' + v.id + '">' + v.nombre + '</option>');
            })
        });
    $('#municipio').html('<option value="--" >Seleccionar...</option>');
    $.post("api/localizacion/consultarmunicipio", data)
        .done(function (data) {
            $.each(data, function (c, v) {
                $('#municipio').append('<option value="' + v.id + '">' + v.nombre + '</option>');
            })
        });
}
//Funcion para lsitar grupos (prueba)
function LisGrupo() {
    //var id = $('#estado option:selected').val();
    var data = JSON.stringify({ ide: parseInt(id) });
    $('#grupo').html('<option value="--" >Seleccionar...</option>');
    $.post("api/listar/grupo", data)
        .done(function (data) {
            $.each(data, function (c, v) {
                $('#grupo').append('<option value="' + v.id + '">' + v.nombre + '</option>');
            })
        });
}
/**
 * Listar Parroquia
 */
function LParroquia() {
    var id = $('#municipio option:selected').val();
    var data = JSON.stringify({ idm: parseInt(id) });
    $('#parroquia').html('<option value="--" >Seleccionar...</option>');
    $.post("api/localizacion/consultarparroquia", data)
        .done(function (data) {
            $.each(data, function (c, v) {
                $('#parroquia').append('<option value="' + v.id + '">' + v.nombre + '</option>');
            })
        });
}

/**
 * Listar Programas
 */
function LProgramas() {
     var table = $('#tbltaquilla').DataTable({
        "paging": false,
        "ordering": false,
        "info": false,
        "searching": false
    });

    var data = JSON.stringify({ id: 3 });

    var table = $('#tblprograma').DataTable({
        "paging": true,
        "ordering": false,
        "info": false,
        "searching": false,
        "pageLength": 6
    });
    $.post("api/listasistema", data)
        .done(function (data) {
            $.each(data, function (c, v) {
                table.row.add([
                    parseInt(c) + 1,
                    v.nombre,
                    '<input type="text" id="triple' + v.oid + '" style="width:60px" maxlength="2">',
                    '<input type="text" id="terminal' + v.oid + '" style="width:60px" maxlength="2">',
                    '<input type="text" id="queda' + v.oid + '" style="width:60px" maxlength="2">',
                    '<input type="text" id="participacion' + v.oid + '" style="width:60px" maxlength="2">'
                ]).draw(false);

            });
        });





}

function CargarCalendario() {
    //Date range picker


    var local = {
        "format": 'YYYY/MM/DD',
        "applyLabel": "Aceptar",
        "cancelLabel": "Cancelar",
        "customRangeLabel": 'Por Rango',
        "daysOfWeek": [
            "Do",
            "Lu",
            "Ma",
            "Mi",
            "Ju",
            "Vi",
            "Sa"
        ],
        "monthNames": [
            "Enero",
            "Febrero",
            "Marzo",
            "Abril",
            "Mayo",
            "Junio",
            "Julio",
            "Agosto",
            "Septiembre",
            "Octubre",
            "Noviembre",
            "Diciembre"
        ],
    };
    //Date range as a button
    $('#daterange-btn').daterangepicker(
        {
            locale: local,
            ranges: {
                'Hoy': [moment(), moment()],
                'Ayer': [moment().subtract(1, 'days'), moment().subtract(1, 'days')],
                'Hace 7 Dias': [moment().subtract(6, 'days'), moment()],
                'Hace 30 Dias': [moment().subtract(29, 'days'), moment()],
                'Este Mes': [moment().startOf('month'), moment().endOf('month')],
                'Mes Pasado': [moment().subtract(1, 'month').startOf('month'), moment().subtract(1, 'month').endOf('month')]
            },
            startDate: moment().subtract(29, 'days'),
            endDate: moment()
        },
        function (start, end) {
            $('#daterange-btn span').html(start.format('YYYY/MM/DD') + ' - ' + end.format('YYYY/MM/DD'));
        }
    );


    $('#fechara').daterangepicker({ locale: local });

    $('#fecha').datepicker({ autoclose: true, format: 'yyyy-mm-dd' });

    $('#fechade').datepicker({ autoclose: true, format: 'yyyy-mm-dd' });
    $('#fechadepositore').datepicker({ autoclose: true, format: 'yyyy-mm-dd' });
    $('#fechadere').datepicker({ autoclose: true, format: 'yyyy-mm-dd' });
}
/**
 * Enviando Archivos
 */
function enviarArchivo() {
    if ($("#archivo").val() == "") {
        $.notify("Debe seleccionar un archivo", "error");
        return false;
    }
    $("#cargando").show();

    var formData = new FormData(document.forms.namedItem("forma"));
    $.ajax({
        url: "subirl",
        type: "post",
        dataType: "html",
        data: formData,
        timeout: 15000,
        cache: false,
        contentType: false,
        processData: false
    })
        .done(function (res) {
            $("#archivo").val("");
            $.notify("Envio de archivos exitosos...", "success");
            $("#cargando").hide();


        }).fail(function (jqXHR, textStatus) {
            $("#archivo").val("");
            if (textStatus === 'timeout') {
                $.notify("Los archivos exceden el limite en tiempo de conexion intente con menos...", "error");
            }
            $("#cargando").hide();
        });

}

/**
 * Listar Reporte de Archivos
 * @param t DataTable
 */
function LstRA() {

    var f = $('#daterange-btn span').html();
    var f_a = f.split("-");
    if (f_a.length < 3) {
        var rfecha = JSON.stringify({
            desde: f_a[0].replace(" ", ""),
            hasta: f_a[1].replace(" ", "")
        });

        var t = $('#reporte').DataTable();
        t.clear().draw();
        $("#cargando").show();
        $.post("api/reportearchivo", rfecha)
            .done(function (data) {
                // Get the column API object
                t.column(1).visible(false);
                t.column(2).visible(false);

                $.each(data, function (c, v) {
                    cantidad = v.cantidad == null ? 0 : v.cantidad;
                    if (v.estatus == "1") estatus = '<span class="label label-success">Procesado</span>';

                    switch (cantidad) {
                        case 0:
                            estatus = '<span class="label label-danger">Sin ventas</span>';
                            break;
                        case 1:
                            estatus = '<span class="label label-success">Procesado</span>';
                            break;
                        default:
                            //estatus = '<span class="label label-info">Pendiente</span>';
                            break;
                    }
                    t.row.add([
                        parseInt(c) + 1,
                        v.oid,
                        v.tabla,
                        v.nombre,
                        cantidad,
                        v.creado,
                        v.procesado,
                        estatus,
                    ]).draw(false);

                });
                $("#cargando").hide();
            });

        var table = $('#reporte').DataTable();
        $('#reporte tbody').on('click', 'tr', function () {
            var data = table.row(this).data();
            //console.log(data);
            VentanaEmergente(data[3], data[1], data[2]);

        });
    } else {
        $.notify("Debe seleccionar un rango", "error");
    }
}

/**
 * Crear Ventanas emergentes con diferentes contenidos
 * @param html
 * @param html
 */
function VentanaEmergente(titulo, id, tbl) {
    var cont = '<table class="table table-bordered table-striped table-hover" \
     id="reporteDetallado" width="100%">\
              <thead>\
                <tr>\
                  <th style="width: 5px">#</th>\
                  <th>Agencia</th>\
                  <th >Venta</th>\
                  <th >Premio</th>\
                  <th >Comisión</th>\
                  <th >Saldo</th>\
                </tr>\
              </thead>\
            </table>';
    $('#ventanaEmergenteTitulo').html(titulo);
    $('#ventanaEmergenteContenido').html(cont);

    $('#ventanaEmergente').modal({ keyboard: false })   // initialized with no keyboard
    $('#ventanaEmergente').modal('show')                // initializes and invokes show
    Reporte('reporteDetallado', id, tbl);
}

/**
 * Crear Reportes
 *
 * @param DataTable
 */
function Reporte(t, oid, tbl) {

    var table = $('#' + t).DataTable({
        "paging": true,
        "ordering": true,
        "info": true,
        "searching": true
    });
    var clave = JSON.stringify({ id: oid, tabla: tbl });
    table.clear().draw();
    $.post("api/reportesaldo", clave)
        .done(function (data) {
            $.each(data, function (c, v) {
                venta = v.ven == null ? 0 : v.ven;
                premio = v.pre == null ? 0 : v.pre;
                comision = v.com == null ? 0 : v.com;
                saldo = v.sal == null ? 0 : v.sal;

                table.row.add([
                    parseInt(c) + 1,
                    v.age,
                    venta,
                    premio,
                    comision,
                    saldo

                ]).draw(false);

            });
        });

}

/**
 * Crear Listado de Notificaciones
 *
 * @param string
 * @param string
 */
function CrearNotificacion(t, msj) {
    $.notify(msj, "info");
    var cant = parseInt($("#hnoti").html()) + 1;
    $("#hnoti").html(cant)
    $("#tnoti").html("Tienes " + cant + " notificaciones")
    $("#cnoti").append('<li> \
                        <a href="#" title="' + msj + '">\
                        <i class="fa fa-check text-aqua"></i>\
                        ' + msj + '\
                        </a>\
                        </li>');


}

/**
 *
 */
function CNErr(t, msj) {
    $.notify(msj, "error");
    var cant = parseInt($("#hnoti").html()) + 1;
    $("#hnoti").html(cant)
    $("#tnoti").html("Tienes " + cant + " notificaciones")
    $("#cnoti").append('<li> \
                        <a href="#" title="' + msj + '">\
                        <i class="fa  fa-bug text-red"></i>\
                        ' + msj + '\
                        </a>\
                        </li>');


}

/**
 * Activar Mensaje de Chat
 *
 * @param string
 * @param string
 * @param Date
 */
function ActivarChat(de, msj, t) {
    $.notify(msj, "info");
}

/**
 * Listar Saldos Generales de Ventas
 *
 * @param Date
 */
function LstSaldo() {

    var f = $('#daterange-btn span').html();
    var suma = 0;
    var aventa = 0;
    var apremio = 0;
    var acomision = 0;

    var f_a = f.split("-");

    if (f_a.length < 3) {
        if ($("#tipo option:selected").val() == "0") {
            $.notify("Debe seleccionar un esquema ", "error");
            return
        }
        var desdeA = f_a[0].replace(" ", "");
        var hastaA = f_a[1].replace(" ", "");

        var rS = $('#reporteSaldos').DataTable();
        var rfecha = JSON.stringify({
            desde: desdeA.replace(/\//g, "-"),
            hasta: hastaA.replace(/\//g, "-"),
            tabla: $("#tipo option:selected").val(),
            sistema: parseInt($("#sistema option:selected").val())
        });
        rS.clear().draw();
        $("#cargando").show();


        $.post("api/reportesaldo", rfecha)
            .done(function (data) {
                $.each(data, function (c, v) {
                    saldo = v.sal == null ? 0 : v.sal;
                    venta = v.ven == null ? 0 : v.ven;
                    premio = v.pre == null ? 0 : v.pre;
                    comision = v.com == null ? 0 : v.com;

                    suma += parseFloat(saldo);
                    aventa += parseFloat(venta);
                    apremio += parseFloat(premio);
                    acomision += parseFloat(comision);

                    rS.row.add([
                        parseInt(c) + 1,
                        v.age,
                        venta,
                        premio,
                        comision,
                        saldo,
                        v.fec,
                    ]).draw(false);

                });

                $("#tfventa").html(aventa.toFixed(2));
                $("#tfpremio").html(apremio.toFixed(2));
                $("#tfcomision").html(acomision.toFixed(2));
                $("#tfsaldo").html(suma.toFixed(2));
                $("#cargando").hide();
            });
    } else {
        $.notify("Debe seleccionar un rango", "error");
    }
}

/**
 *
 */
function LSistema() {
    var tipo = 0;
    console.log($("#tipo option:selected").val());
    switch ($("#tipo option:selected").val()) {
        case "loteria":
            tipo = 0
            break;
        case "parley":
            tipo = 1;
            break;
        case "figura":
            tipo = 2;
            break;
        case "truco":
            tipo = 3;
            break;
        case "pescalo":
            tipo = 4;
            break;

        case "todos":
            tipo = 5;
            break;
        default:
            $("#sistema").html("<option value='--'>------------</option>");
            return;
    }

    var data = JSON.stringify(
        {
            id: tipo
        }
    );

    $.post("api/listasistema", data)
        .done(function (data) {
            $("#sistema").html("");
            $.each(data, function (c, v) {
                $("#sistema").append("<option value='" + v.oid + "'>\
               " + v.nombre + "</option>");

            });
            $("#sistema").append("<option value=99>Todos</option>");
        });

}

/**
 *
 */
function LstSaldoGPS() {

    var f = $('#daterange-btn span').html();
    var suma = 0;

    var f_a = f.split("-");

    if (f_a.length < 3) {
        if ($("#tipo option:selected").val() == "--") {
            $.notify("Debe seleccionar un esquema ", "error");
            return
        }

        var rfecha = JSON.stringify({
            id: parseInt($("#tipo option:selected").val()),
            desde: f_a[0].replace(" ", ""),
            hasta: f_a[1].replace(" ", "")
        });
        url = "api/reportesaldogeneral";
        if (parseInt($("#tipo option:selected").val()) == 2) url = "api/balancegeneral";

        $.post(url, rfecha)
            .done(function (data) {
                switch (parseInt($("#tipo option:selected").val())) {
                    case 0:
                        PLoteria(data);
                        break;
                    case 1:
                        PParley(data);
                        break;
                    case 2:
                        PTotales(f_a[0].replace(" ", ""), f_a[1].replace(" ", ""), data);
                        break;
                    default:
                        break;
                }

            });
    } else {
        $.notify("Debe seleccionar un rango", "error");
    }
}

/**
 *
 */
function PLoteria(data) {
    $("#divTabla").html('\
        <table class="table table-bordered" id="reporteSaldosGeneral" width="100%">\
            <thead>\
            <tr>\
                <th style="width: 5px">#</th>\
                <th>Fecha</th>\
                <th>Morpheus</th>\
                <th>Pos 1</th>\
                <th>Pos 2</th>\
                <th>Pos 3</th>\
                <th>Maticlo</th>\
                <th>Total</th>\
            </tr>\
            </thead>\
        </table>');
    var rS = $('#reporteSaldosGeneral').DataTable();
    rS.clear().draw();
    var i = 0;
    $.each(data, function (c, v) {
        var fila = {};
        var suma = 0;
        $.each(v, function (cl, va) {
            saldo = va.saldo == null ? 0 : va.saldo;
            fila[va.sistema] = saldo;
            suma += saldo;
        });

        morpheus = fila[1] == null ? 0 : fila[1]
        pos1 = fila[2] == null ? 0 : fila[2]
        pos2 = fila[3] == null ? 0 : fila[3]
        pos3 = fila[4] == null ? 0 : fila[4]
        maticlo = fila[5] == null ? 0 : fila[5]
        i++;
        rS.row.add([
            i,
            c,
            morpheus,
            pos1,
            pos2,
            pos3,
            maticlo,
            suma
        ]).draw(false);


    });
}

/**
 *
 */
function PParley(data) {
    $("#divTabla").html('\
        <table class="table table-bordered" id="reporteSaldosGeneral" width="100%">\
            <thead>\
            <tr>\
                <th style="width: 5px">#</th>\
                <th>Fecha</th>\
                <th>Ilbanquero</th>\
                <th>CyberParley</th>\
                <th>Sport17</th>\
                <th>Total</th>\
            </tr>\
            </thead>\
        </table>');
    var rS = $('#reporteSaldosGeneral').DataTable();
    rS.clear().draw();
    var i = 0;
    $.each(data, function (c, v) {
        var fila = {};
        var suma = 0;
        $.each(v, function (cl, va) {
            saldo = va.saldo == null ? 0 : va.saldo;
            fila[va.sistema] = saldo;
            suma += saldo;
        });

        ilbanquero = fila[6] == null ? 0 : fila[6]
        cyberparley = fila[7] == null ? 0 : fila[7]
        sport = fila[8] == null ? 0 : fila[8]

        i++;
        rS.row.add([
            i,
            c,
            ilbanquero,
            cyberparley,
            sport,
            suma
        ]).draw(false);


    });
}

/**
 *
 */
function psFila(fila, buscar) {
    var pos = 0;

    $.each(fila, function (c, v) {
        if (buscar == v) {
            pos = c;
        }
    });
    return pos;
}

/**
 * Listar Totales de los saldos para estado de Cuenta
 *
 */
function PTotales(desde, hasta, data) {
    rS = CBTotales();
    var fila = RecorreFechas(desde, hasta, rS);
    var acumuladorSaldosFecha = {};
    var i = 0;
    var total = 0;
    $.each(data, function (c, v) {
        //var fila = {};
        var suma = 0;
        var pos = [];
        pos[0] = 2; //Loteria
        pos[1] = 3;  //Parley

        $.each(v, function (cl, va) {

            //console.log(va);
            saldo = va.saldo == null ? 0 : va.saldo;

            debe = va.debe == null ? 0 : va.debe;
            haber = va.debe == null ? 0 : va.haber;
            fil = psFila(fila, cl) - 1;

            //console.log("Fila: " + fil + " Columna: "+ pos[c] + " Fecha: " + cl);
            rS.cell(fil, pos[c]).data(saldo).draw();

            rS.cell(fil, 5).data(debe).draw();
            rS.cell(fil, 6).data(haber).draw();
            total = parseFloat(rS.cell(fil, 2).data());
            total += parseFloat(rS.cell(fil, 3).data());
            rS.cell(fil, 4).data(total).draw();

        });
        total = 0;

    });


}

/**
 *
 */
function CBTotales() {
    $("#divTabla").html('\
    <table class="table table-bordered" id="reporteSaldosGeneral" width="100%">\
        <thead>\
        <tr>\
            <th style="width: 5px">#</th>\
            <th>Fecha</th>\
            <th>Loteria</th>\
            <th>Parley</th>\
            <th>Total</th>\
            <th>Entregado</th>\
            <th>Recibido</th>\
            <th>Balance</th>\
        </tr>\
        </thead>\
    </table>');
    var rS = $('#reporteSaldosGeneral').DataTable();
    rS.clear().draw();
    //rS.cells(2,3).data("Hola").draw();
    return rS;
}

/**
 *
 */
function GC(tipo) {
    if ($("#fecha").val() == "") {
        $.notify("Debe seleccionar la fecha", "error");
        return
    }
    $("#cargando").show();
    $("#divReporte").html(tableGC());
    $("#reporte").DataTable(opciones);
    var t = $("#reporte").DataTable();
    fecha = $("#fecha").val();
    fecha = fecha.replace(/-/g, "/");
    fecha = OperarFecha(fecha, -1);
    $("#lblFechade").html(fecha);
    var data = JSON.stringify({
        fecha: fecha
    });
    url = evalTipo();
    $.post(url, data)
        .done(function (data) {
            t.clear().draw();
            var i = 0;

            $.each(data, function (c, v) {

                vienen = v.vienen == null ? 0 : v.vienen;
                saldo = v.saldo == null ? 0 : v.saldo;
                ingreso = v.ingreso == null ? 0 : v.ingreso;
                egreso = v.egreso == null ? 0 : v.egreso;
                prestamo = v.prestamo == null ? 0 : v.prestamo;
                entregado = v.entregado == null ? 0 : v.entregado;
                recibido = v.recibido == null ? 0 : v.recibido;
                cuota = v.cuota == null ? 0 : v.cuota;
                movimiento = (parseFloat(egreso) + parseFloat(cuota)) - (parseFloat(ingreso) + parseFloat(prestamo));

                x = parseFloat(entregado) - parseFloat(recibido);
                //console.log("SALDO: " + v.saldo + " X: " + x + " MOVIMIENTO : " + movimiento);
                total = vienen + parseFloat(saldo) + movimiento + x;

                if ($("#txtSeleccion").val() == "0"){
                    accion = btnAccion(c, v.observacion, total);
                    nombre =  v.observacion;
                }else{
                    accion = btnAccion(v.oid, v.agencia, total);
                    nombre = v.agencia;
                }

                i++;
                if(i == 1){

                  if (v.estatus != null) {

                      if (v.estatus == 1) {
                          accion = "";
                          t.column(0).visible(false); //Ocultar la columna 0
                      }
                  }
                }


                switch (parseInt(tipo)) {
                    case 0:
                        //console.log('CERO... ');
                        t.row.add([
                            accion,
                            nombre,
                            vienen.toFixed(2),
                            saldo.toFixed(2),
                            movimiento.toFixed(2),
                            x.toFixed(2),
                            total.toFixed(2)
                        ]).draw();
                        break;
                    case 1:
                        //console.log('UNO... ');
                        if (total >= 0) {
                            t.row.add([
                                accion,
                                nombre,
                                vienen.toFixed(2),
                                saldo.toFixed(2),
                                movimiento.toFixed(2),
                                x.toFixed(2),
                                total.toFixed(2)
                            ]).draw();
                        }
                        break;
                    default:
                        //console.log('DOS... ');
                        if (total < 0) {
                            t.row.add([
                                accion,
                                nombre,
                                vienen.toFixed(2),
                                saldo.toFixed(2),
                                movimiento.toFixed(2),
                                x.toFixed(2),
                                total.toFixed(2)
                            ]).draw();
                        }
                        break;
                }


            });
            $("#cargando").hide();

        })
}


function evalTipo(){
    var url = "";
    switch ($("#txtSeleccion").val()) {
        case "0":
            url = "api/balance/cobrosypagosgrupo";
            break;
        case "1":
            url =  "api/balance/cobrosypagoscolector";
            break;
        case "2":
            url =  "api/balance/cobrosypagos";
            break;
        default:
            url = "api/balance/cobrosypagosgrupo";
            break;
    }
    return url;
}

/**
 *
 */
function tableGC() {
    s = '<table class="table table-bordered" cellspacing="0" id="reporte" width="100%">\
            <thead>\
              <tr>\
                <th style="width: 60px">#</th>\
                <th>Grupo</th>\
                <th>S. Ant.</th>\
                <th>S. Día</th>\
                <th>Movimiento</th>\
                <th>+ E - R</th>\
                <th>Total</th>\
              </tr>\
            </thead>\
          </table>';
    return s;

}

/**
 *
 */
function btnAccion(valor, texto, monto) {
    s = '<div class="btn-group">\
        <button type="button" class="btn btn-success">\
        <span class="fa fa-cogs"></span></button>\
        <button type="button" class="btn btn-success dropdown-toggle" \
        data-toggle="dropdown" aria-expanded="false">\
        <span class="caret"></span>\
        <span class="sr-only">Toggle Dropdown</span>\
        </button>\
        <ul class="dropdown-menu" role="menu">\
            <li><a href="#" onclick="mdlE(\'mdlMovimiento\',\'\' , \'' + valor + '\', \'' + monto + '\', \'' + texto + '\')">Registrar Movimiento</a></li>\
            <li class="divider"></li>\
            <li><a href="#" onclick="mdlE(\'mdlER\',\'er\', \'' + valor + '\', \'' + monto + '\', \'' + texto + '\')">Registrar +E -R</a></li>\
            <li><a href="#" onclick="mdlE(\'mdlPre\',\'pre\', \'' + valor + '\', \'' + monto + '\', \'' + texto + '\')">Registrar Prestamos </a></li>\
            <li><a href="#" onclick="mdlE(\'mdlEC\',\'ec\', \'' + valor + '\', \'' + monto + '\', \'' + texto + '\');">Estado de Cuenta </a></li>\
        </ul>\
    </div>';
    return s
}


function GCD() {
    fecha = $("#fecha").val();
    fecha = fecha.replace(/-/g, "/");
    fecha = OperarFecha(fecha, -1);
    var data = JSON.stringify({
        fecha: fecha,
        cierre: 1
    });
    console.log(data);

    $.post("api/balance/cierrediario", data)
        .done(function (data) {
            $.notify("Proceso exitos: Se han generado todos los eventos del día siguiente...", "success");
    });

}
/**
 */
function mdlE(id, cod, valor, monto, texto) {

    $('#cod' + cod).html(valor);
    $('#cod' + cod + 't').html(texto);
    $('#' + id).modal('show');
    $('#montoer').val(monto);

    $("#fechadepositore").val($("#lblFechade").html());
    $("#fechadere").val($("#fecha").val());

    var msj = 'Saldo a cero (0)';
    if (monto > 0) {
        msj = 'Saldo a cero (0)';
    }
    $('#descripcioner').val(msj);
    $('#divTablaec').html('');
}

/**
 *
 */
function DP() {
    if ($("#fecha").val() == "") {
        $.notify("Debe seleccionar la fecha", "error");
        return
    }
    $('#mdlDP').modal('show');
    tabla = $("#rptDeposito").DataTable();
    var data = JSON.stringify({ fdeposito: $("#fecha").val() });
    $.post("api/movimiento/listardeposito", data)
        .done(function (data) {
            tabla.clear().draw();
            $.each(data, function (c, v) {
                console.log(v);
                banco = v.banco == null ? 0 : v.banco;
                accion = btnADep(v.oid);
                tabla.row.add([
                    v.oid,
                    accion,
                    v.agencia,
                    v.banconombre,
                    v.voucher,
                    v.monto
                ]).draw();
            });
            tabla.column(0).visible(false);
        });
}

/**
 *
 */
function btnADep(oid) {
    s = '<div class="btn-group">\
        <button type="button" class="btn btn-success" onclick="ADP(\'' + oid + '\',\'\')">\
        <span class="fa fa-check-circle"></span></button>\
        <button type="button" class="btn btn-danger" onclick="ODP(\'' + oid + '\')">\
        <span class="fa fa-times-circle"></span></button></div>';
    return s
}

/**
 *
 */
function ADP(oid, obse) {
    fecha = new Date();
    mes = fecha.getMonth() + 1;
    fecha = fecha.getFullYear() + "-" + mes + "-" + fecha.getDate();
    estatus = 1;

    if (obse != "") {
        obse = $('#txtObservacionDep').val();
        $('#msgbox').modal('hide');
        estatus = 2;
    } else {
        obse = "Aprobado";
    }
    var data = JSON.stringify({
        oid: parseInt(oid),
        fecha: fecha,
        formadepago: 0,
        estatus: estatus,
        observacion: obse
    })

    $('#mdlDP').modal('hide');


    $.post("api/movimiento/actualizarer", data)
        .done(function (data) {
            $.notify("El pago ha sido aprobado con exito...", "success");

        });

}

/**
 * Observacion de las aceptaciones o rechazo de los depositos
 * @param int
 */
function ODP(oid) {

    $('#mdlDP').modal('hide');
    $('#msgbox-cuerpo').html('<div class="form-group">\
                  <label>Observaciones</label>\
                  <textarea id="txtObservacionDep" class="form-control" rows="3" placeholder="Observaciones ..."></textarea>\
                </div>');
    $('#msgbox-pie').html('<button type="button" class="btn btn-success" onclick="ADP(\'' + oid + '\',\'txtObservacionDep\')">Procesar</button>\
            <button type="button" class="btn btn-default" data-dismiss="modal">Cerrar</button>');
    $('#msgbox').modal('show');
}

/**
 *
 */
function EC() {
    $('#cagandoec').show();
    //$("#cagandoec").hide();
    var f = $('#fecharangoec option:selected').val();

    var suma = 0;


    var f_a = f.split("-");
    if (f_a.length < 3) {
        desdeA = f_a[0].replace(" ", "");
        hastaA = f_a[1].replace(" ", "");
        if (f_a[0] == "0") {
            desdeA = moment().format('YYYY/MM/') + '01';
            hastaA = moment().format('YYYY/MM/DD');
        }

        var rfecha = JSON.stringify({
            agencia: $('#codec').html(),
            desde: desdeA.replace(/\//g, "-"),
            hasta: hastaA.replace(/\//g, "-")
        });
        /*
        var rfecha =JSON.stringify({
            agencia :  $('#codec').html(),
            desde:f_a[0].replace(" ", ""),
            hasta:f_a[1].replace(" ", "")
        });
        */
        url = "api/balance/cobrosypagos";
        $.post(url, rfecha)
            .done(function (data) {

                PTotalesDetalles(desdeA, hastaA, data);

            });
    } else {
        $.notify("Debe seleccionar un rango", "error");
    }
}

/**
 *
 */
function EstadoCuenta() {
    $("#divTablaec").html('\
    <table class="table table-bordered" id="reporteSaldosGeneral" width="100%">\
        <thead>\
        <tr>\
            <th>Fecha</th>\
            <th>S.Ant</th>\
            <th>S.Día</th>\
            <th>Moviemiento</th>\
            <th>Recibido</th>\
            <th>Entregado</th>\
            <th>Total</th>\
        </tr>\
        </thead>\
    </table>');
    var rS = $('#reporteSaldosGeneral').DataTable({
        scrollY: "200px",
        scrollCollapse: true,
        paging: false,
        searching: true,
        order: [[0, "desc"]]
    });
    rS.clear().draw();
    return rS;
}

/**
 *
 */
function TablaEstadoCuenta(fila, buscar) {
    var pos = 0;

    $.each(fila, function (c, v) {
        if (buscar == v) {
            pos = c;
        }
    });
    return pos;
}

/**
 * Listar Totales de los saldos para estado de Cuenta
 *
 */

function PTotalesDetalles(desde, hasta, data) {
    rS = EstadoCuenta();
    var fila = RecorreFechas(desde, hasta, rS);
    var acumuladorSaldosFecha = {};
    var i = 0;
    var total = 0;
    if (data == null) {
        $.notify("No se encontrarón registros", "error")
        return
    }
    $.each(data, function (c, v) {
        sAnt = v.vienen == null ? 0 : v.vienen;
        ingreso = v.ingreso == null ? 0 : v.ingreso;
        egreso = v.egreso == null ? 0 : v.egreso;
        prestamo = v.prestamo == null ? 0 : v.prestamo;
        entregado = v.entregado == null ? 0 : v.entregado;
        recibido = v.recibido == null ? 0 : v.recibido;
        cuota = v.cuota == null ? 0 : v.cuota;
        movimiento = (parseFloat(egreso) + parseFloat(cuota)) - (parseFloat(ingreso) + parseFloat(prestamo));
        saldo = v.saldo == null ? 0 : v.saldo;
        x = parseFloat(entregado) - parseFloat(recibido);
        total = parseFloat(v.saldo) + movimiento + x;

        fil = psFila(fila, v.fecha);
        rS.cell(fil, 1).data(sAnt.toFixed(2)).draw();

        rS.cell(fil, 2).data(saldo.toFixed(2)).draw();
        rS.cell(fil, 3).data(movimiento.toFixed(2)).draw();
        rS.cell(fil, 4).data(entregado.toFixed(2)).draw();
        rS.cell(fil, 5).data(recibido.toFixed(2)).draw();

        rS.cell(fil, 6).data(total.toFixed(2)).draw();


    });


}

/**
 * @param Date | UNIX
 * @param Date | UNIX
 * @param DataTable
 */
function RecorreFechas(desde, hasta, rS) {

    fauxd = desde.split("/");
    fauxh = hasta.split("/");

    danio = parseInt(fauxd[0]);
    dmes = parseInt(fauxd[1]);
    ddia = parseInt(fauxd[2]);

    hanio = parseInt(fauxh[0]);
    hmes = parseInt(fauxh[1]);
    hdia = parseInt(fauxh[2]);
    var fila = {};
    var count = 0;
    for (h = danio; h <= hanio; h++) {
        for (i = dmes; i <= hmes; i++) {
            mdmes = new Date(h, i, 0).getDate();
            for (j = ddia; j <= mdmes; j++) {
                dia = j;
                if ((String(j)).length == 1) dia = '0' + j;
                mes = i;
                if ((String(i)).length == 1) mes = '0' + i;
                fecha = h + "-" + mes + "-" + dia;
                fila[count] = fecha;
                rS.row.add([
                    fecha,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                ]).draw(false);
                count++;
                //console.log(danio + "-" + i + "-" + j);
                if (hanio == h && hmes == i && hdia == j) break;
            }
            ddia = 1;
        }
    }
    return fila;
}

/**
 *
 */
function RegistrarER() {
    monto = parseFloat($("#montoer").val());
    if ( monto < 0){
        monto = parseFloat($("#montoer").val()) *- 1;
    }
    grupo = 0;
    agencia = 0;
    sel = parseInt($("#txtSeleccion").val());
    codigo = parseInt($("#coder").html());
    if(sel == 0 ){
      grupo = codigo;
    }else if(sel == 2){
      agencia = codigo;
    }
    var EntregadoRecibido = JSON.stringify({
        oid : agencia,
        subgrupo : 0,
        colector : 0,
        grupo: grupo,
        estatus: 1,
        fecha: $("#fechadere").val(),
        fechaaprobado: $("#lblFechade").html(),
        fechaoperacion: $("#fecha").val(),
        deposito: $("#fechadepositore").val(),
        forma: parseInt($("#tipoer option:selected").val()), //0 Entregado: DEBE 1 Recibido:HABER
        banco: parseInt($("#cuentaer option:selected").val()),
        monto: monto,
        voucher: $("#voucer").val(),
        observacion: $("#descripcioner").val()
    });
    console.log(EntregadoRecibido);
    url = "api/balance/registrarpago";
    $.post(url, EntregadoRecibido)
        .done(function (data) {
            $('#mdlER').modal('hide');
            $.notify("El registro ha sido exitoso, si desea verlo en pantalla presione F5.", "success");ss
            $("#voucer").val('');
            GC(0);
        });

}


/**
 * Funcion que devuelve la fecha actual y la fecha modificada n dias
 * Tiene que recibir el numero de dias en positivo o negativo para sumar o
 * restar a la fecha actual.
 * Ejemplo:
 *  mostrarFecha('YYYY/MM/DD', -10) => restara 10 dias a la fecha actual
 *  mostrarFecha('YYYY/MM/DD', 30) => añadira 30 dias a la fecha actual
 *
 * @param date
 * @param int
 */
function OperarFecha(fecha, dias) {
    milisegundos = parseInt(35 * 24 * 60 * 60 * 1000);
    fecha = new Date(fecha);
    //Obtenemos los milisegundos desde media noche del 1/1/1970
    tiempo = fecha.getTime();
    //Calculamos los milisegundos sobre la fecha que hay que sumar o restar...
    milisegundos = parseInt(dias * 24 * 60 * 60 * 1000);
    //Modificamos la fecha actual
    total = fecha.setTime(tiempo + milisegundos);
    day = fecha.getDate();
    month = fecha.getMonth() + 1;
    year = fecha.getFullYear();

    return year + "-" + month + "-" + day;
}

/**
 * Recibo de Pagos
 */
function RPago() {
    dep = $("#fechade").val();
    mon = $("#monto").val();
    vouc = $("#numoperacion").val();
    if (dep == "") {
        $.notify("Debe introducir una fecha ", "warn");
        return false;
    }
    if (mon == "") {
        $.notify("Debe introducir un monto ", "warn");
        return false;
    }

    var Pago = JSON.stringify({
        comercializadora: 1,
        grupo: 0,
        subgrupo: 0,
        colector: 0,
        agenciacod: 0,
        voucher: vouc,
        fecha: dep,
        cuentadebe: parseInt($("#cuentadebe").val()),
        tipodebe: parseInt($("#tipodebe").val()),
        cuentahaber: parseInt($("#cuentahaber").val()),
        tipohaber: parseInt($("#tipohaber").val()),
        monto: parseFloat(mon),
        observacion: $("#descripcion").val()
    });

    $.post("api/movimiento/registrar", Pago)
        .done(function (data) {
            $("#cuentadebe option:selected").val('--');
            $("#cuentahaber option:selected").val('--');
            $("#tipodebe option:selected").val('--');
            $("#tipohaber option:selected").val('--');
            $("#monto").val('');
            //$("#fechade").val('');
            $("#descripcion").val('');
            $("#numoperacion").val('');
            $.notify("Se ha registrado el movimiento", "success");
            LPago();
        });

}

/**
 * Registrar Movimiento Individual
 */
function RPagoMI() {
    dep = $("#fechade").val();
    mon = $("#monto").val();
    vouc = $("#numoperacion").val();
    if (dep == "") {
        $.notify("Debe introducir una fecha ", "warn");
        return false;
    }
    if (mon == "") {
        $.notify("Debe introducir un monto ", "warn");
        return false;
    }

    var Pago = JSON.stringify({
        comercializadora: 0,
        grupo: parseInt($("#cod").html()),
        subgrupo: 0,
        colector: 0,
        agenciacod: 0,
        voucher: vouc,
        fecha: dep,
        cuentadebe: parseInt($("#cuentadebe").val()),
        tipodebe: parseInt($("#tipodebe").val()),
        cuentahaber: parseInt($("#cuentahaber").val()),
        tipohaber: parseInt($("#tipohaber").val()),
        monto: parseFloat(mon),
        observacion: $("#descripcion").val()
    });

    $.post("api/movimiento/registrar", Pago)
        .done(function (data) {
            $("#cuentadebe option:selected").val('--');
            $("#cuentahaber option:selected").val('--');
            $("#tipodebe option:selected").val('--');
            $("#tipohaber option:selected").val('--');
            $("#monto").val('');
            $("#descripcion").val('');
            $("#numoperacion").val('');
            $.notify("Se ha registrado el movimiento", "success");
            LPago();
        });

}

/**
 * Listar Recibos de Pagos Movimiento (Entregado y Recibidos)
 */

function LPago() {
    dep = $("#fechade").val();
    var Pago = JSON.stringify({
        fecha: dep
    });
    $("#divReporte").show();
    var rS = $('#lstReporte').DataTable();
    rS.clear().draw();

    $.post("api/movimiento/listar", Pago)
        .done(function (data) {
            $.each(data, function (c, v) {
                token = v.token;
                rS.row.add([
                    v.voucher,
                    v.cuentadeben,
                    STipo(v.tipodebe),
                    v.cuentahabern,
                    STipo(v.tipohaber),
                    v.observacion,
                    parseFloat(v.monto).toFixed(2)
                ]).draw();
            });
        });
}

/**
 *
 */
function STipo(id) {
    switch (id) {
        case 1:
            return 'Deposito';
            break;
        case 2:
            return 'Transferencia';
            break;
        case 3:
            return 'Cheque';
            break;
        case 4:
            return 'Ingreso';
            break;
        case 5:
            return 'Egreso';
            break;
        default:
            return 'Deposito';
            break;
    }
}


function LimpiarGrupo() {
    $("#nombregrupo").val('');
    $("#fecha").val('');
    $("#cuenta").val('');
    $("#terminal").val('');
    $("#triple").val('');
    $("#queda").val('');
    $("#participacion").val('');
    $("#observacion").val('');

    $("#casa").val('');
    $("#direccion").val('');
    $("#telefono").val('');
    $("#celular").val('');

    $("#usuario").val('');
    $("#correo").val('');
    $("#clave").val('');
    $("#rclave").val('');
    $("#respuesta").val('');
    LEstado();

}

function ValidarRegistro() {

    if ($("#nombregrupo").val() == "") {
        $.notify("Debe ingresar un nombre para el grupo", "error");
        return false;
    }
    if ($("#fecha").val() == "") {
        $.notify("Debe ingresar una fecha", "error");
        return false;
    }
    if ($("#parroquia").val() == "--") {
        $.notify("Debe seleccionar una localización", "error");
        return false;
    }

    if ($("#usuario").val() == "") {
        $.notify("Debe ingresar un nombre de usuario", "error");
        return false;
    }

    if ($("#clave").val() != "") {
        if ($("#clave").val() != $("#rclave").val()) {
            $.notify("La contraseñas no son identica", "error");
            return false;
        }
    } else {
        $.notify("La contraseñas está en blanco", "error");
        return false;
    }


    return true
}

function RegistrarGrupo() {
    if (ValidarRegistro() == false) return;

    var Localizacion = {
        idp: parseInt($("#parroquia").val()),
        casa: $("#casa").val(),
        direccion: $("#direccion").val(),
        telefono: $("#telefono").val(),
        celular: $("#celular").val(),
        tipo: parseInt($("#tipo").val())
    };
    var Seguridad = {
        usuario: $("#usuario").val(),
        correo: $("#correo").val(),
        clave: $("#clave").val(),
        rclave: $("#rclave").val(),
        pregunta: parseInt($("#pregunta").val()),
        respuesta: $("#respuesta").val()
    };

    var Grupo = JSON.stringify({
        nombre: $("#nombregrupo").val(),
        fecha: $("#fecha").val(),
        cuenta: $("#cuenta").val(),
        terminal: parseFloat($("#terminal").val()),
        triple: parseFloat($("#triple").val()),
        queda: parseFloat($("#queda").val()),
        participacion: parseFloat($("#participacion").val()),
        frecuencia: parseInt($("#frecuencia").val()),
        negociacion: parseInt($("#negociacion").val()),
        observacion: $("#observacion").val(),
        localizacion: Localizacion,
        seguridad: Seguridad
    });

    $("#cargando").show();
    $.post("api/registro/grupo", Grupo)
        .done(function (data) {
            if (data.tipo != 2) {
                $.notify("Envio de archivos exitosos...", "success");
            } else {
                $.notify(data.msj, "error");
            }
            LimpiarGrupo();
            $('#tabgrupo a:first').tab('show') // Select first tab
            $("#cargando").hide();
        });
}

function ValidarPar() {
    if (parseFloat($("#participacion").val()) > 0) {
        $("#queda").val('0');
    }
}
function ValidarQueda() {
    if (parseFloat($("#queda").val()) > 0) {
        $("#participacion").val('0');
    }
}

/**
 * Agregar codigos de cajas
 */
function agregarCajas() {
    var table = $('#tbltaquilla').DataTable();
    var taquilla = $("#taquilla").val();
    var vendedor = $("#vendedor").val();
    table.row.add([taquilla, vendedor]).draw(false);
    $("#taquilla").val("");
    $("#vendedor").val("");

}

function ListarGrupos(grupo){
    var s = '<div class="col-md-4">\
    <div class="box box-widget widget-user-2">\
            <div class="widget-user-header bg-aqua">\
              <div class="widget-user-image">\
                <img class="img-circle" src="dist/img/user2-160x160.jpg" alt="User Avatar">\
              </div>\
              <h3 class="widget-user-username">' + grupo + '</h3>\
              <h5 class="widget-user-desc">Lead Developer</h5>\
            </div>\
            <div class="box-footer no-padding">\
              <ul class="nav nav-stacked">\
                <li><a href="#">Projects <span class="pull-right badge bg-blue">31</span></a></li>\
                <li><a href="#">Tasks <span class="pull-right badge bg-aqua">5</span></a></li>\
                <li><a href="#">Completed Projects <span class="pull-right badge bg-green">12</span></a></li>\
                <li><a href="#">Followers <span class="pull-right badge bg-red">842</span></a></li>\
              </ul>\
            </div>\
          </div>\
        </div>';
    return s;
}



/**
 * Seleccionar Grupo
 */
function SeGrupo(){
    $("#txtSeleccion").val("0");
    $("#btnSeleccion").html('Grupo&nbsp;&nbsp;<span class="fa fa-caret-down"></span>');
}

/**
 * Seleccionar el Colector
 */
function SeColector(){
    $("#txtSeleccion").val("1");
    $("#btnSeleccion").html('Colector&nbsp;&nbsp;<span class="fa fa-caret-down"></span>');
}


/**
 * Seleccionar la Agencia
 */
function SeAgencia(){
    $("#txtSeleccion").val("2");
    $("#btnSeleccion").html('Agencia&nbsp;&nbsp;<span class="fa fa-caret-down"></span>');
}
