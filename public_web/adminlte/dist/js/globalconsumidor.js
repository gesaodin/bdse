var _S = ['Domingo', 'Lunes', 'Martes', 'Miércoles', 'Jueves', 'Viernes', 'Sábado'];
var _SAV = ['Dom', 'Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb'];

var opciones = {
        "paging":   true,
        "ordering": true,
        "info":     true,
        "searching": false
    };
var opcionesfalso = {
        "paging":   false,
        "ordering": false,
        "info":     false,
        "searching": false
    };
//$('#reporte').DataTable(opciones);
$('#reporteSaldos').DataTable(opciones);
$('#reporteSaldosSistemas').DataTable(opcionesfalso);
$('#reporteSaldosSistemasParley').DataTable(opcionesfalso);
 idleTime = 0;
$(function(){
    /**
     * Declaración de variables globales
     */
    //var rS = $('#reporteSaldos').DataTable();



    //Increment the idle time counter every minute.

    var idleInterval = setInterval("timerIncrement()", 60000); // 1 minute
    //Zero the idle timer on mouse movement.
    $(this).mousemove(function (e) {
        //console.log('Entrando...');
        idleTime = 0;
    });
    $(this).keypress(function (e) {
        idleTime = 0;
    });

    if ($('#bancrp').val() != undefined) LCuentaM();


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


    CargarCalendario();
});


/**
 * Listar Cuentas para Movimientos
 */
function LCuentaM(id) {
    $('#cuentadebe').html('<option value="--" >Seleccionar...</option>');
    var data = JSON.stringify({operacion:0});
    $.post("api/movimiento/listarcuentas", data)
        .done(function (data) {
            $.each(data, function (c, v) {
                $('#cuentadebe').append('<option value="' + v.oid + '">' + v.oid + ' | ' + v.nombre + '</option>');
            })
        });
}


function CargarCalendario(){


    $('#fecha').datepicker({autoclose: true, format: 'yyyy-mm-dd'});
    $('#fechade').datepicker({autoclose: true, format: 'yyyy-mm-dd' });
    $('#fechadere').datepicker({autoclose: true, format: 'yyyy-mm-dd' });
}



function timerIncrement() {
    idleTime++;
    if(idleTime > 4)$(location).attr('href', 'logout');
}


//Listar pagos pendientes del registro de pago.
function LPago(){
    $("#divReporte").show();
    var rS = $('#lstReporte').DataTable();
    var Pago = JSON.stringify({
        agencia: $("#agencia").html()
    });
    rS.clear().draw();
    $.post("api/balance/listarpagos", Pago)
    .done(function (data){
        var i = 1;
        $.each(data, function(c, v) {

            estatus = v.estatus == null?0:v.estatus;
            fechaaprobado = v.fechaaprobado == "null"?'':v.fechaaprobado;
            observacion = v.estatus == null?0:v.estatus;
            switch (estatus) {
                    case 0:
                        estatus = '<span class="label label-warning">Pendiente</span>';
                        break;
                    case 1:
                        estatus = '<span class="label label-success">Procesado</span>';
                        break;
                    case 2:
                        estatus = '<span class="label label-danger">Rechazado</span>';
                        fechaaprobado = v.observacion;
                        break;
                    default:
                        //estatus = '<span class="label label-info">Pendiente</span>';
                        break;
            }

            rS.row.add([
                i++,
                v.banconombre,
                v.fecha,
                fechaaprobado,
                v.voucher,
                parseFloat(v.monto).toFixed(2),
                estatus
            ]).draw();
        });
    });


}

//Registrar Pagos
function RPago(){
    dep = $("#fechade").val();
    mon = $("#mont").val();
    vouc = $("#vouc").val();
    if(dep == ""){
        $.notify("Debe introducir una fecha ", "warn");
        return false;
    }
    if(mon == ""){
        $.notify("Debe introducir un monto ", "warn");
        return false;
    }

    var Pago = JSON.stringify({
        agencia: $("#agencia").html(),
        voucher: vouc,
        deposito: dep,
        banco : parseInt($("#banc").val()),
        formadepago: parseInt($("#form").val()),
        monto: parseFloat(mon),
        observacion: $("#obse").val()
    });
    $.post("api/balance/registrarpago", Pago)
    .done(function (data){
        $.notify("Proceso exitoso ", "success");
        LFPago();
        LPago();
    });
}


function RP(){
    //$("#mdlRP").modal('show');
    $(location).attr('href', 'registropago');
}
/**
 * Limpiar el formulario de pago
 */
function LFPago(){
    $("#fechade").val("");
    $("#vouc").val("");
    $("#obse").val("");
    $("#mont").val("");
}

function LstSaldoGPS(){

    var f = $('#fecharango option:selected').val();
    var suma = 0;

    var f_a = f.split("-");
    if (f_a.length < 3){
        desdeA = f_a[0].replace(" ", "");
        hastaA = f_a[1].replace(" ", "");
        if (f_a[0] == "0"){
            desdeA = moment().format('YYYY/MM/') + '01';
            hastaA = moment().format('YYYY/MM/DD');
        }

        var rfecha =JSON.stringify({
            agencia : $("#agencia").html(),
            desde: desdeA.replace(/\//g, "-"),
            hasta: hastaA.replace(/\//g, "-")
        });
        url = "api/balance/cobrosypagos";


        $.post(url,rfecha)
        .done(function(data){

            PTotales(desdeA, hastaA, data);

        });



    }else{
        $.notify("Debe seleccionar un rango", "error");
    }
}

function psFila(fila, buscar){
    var pos = 0;

    $.each(fila, function(c, v){
        if (buscar == v){
            pos = c;
        }
    });
    return pos;
}

/**
 * Listar Totales de los saldos para estado de Cuenta
 *
 */

function PTotales(desde, hasta, data){
    rS = CBTotales();
    var fila = RecorreFechas(desde, hasta, rS);
    var acumuladorSaldosFecha = {};
    var i = 0;
    var total = 0;
     if (data == null){

        $.notify("No se encontrarón registros", "error")
        return
    }
    $.each(data, function(c, v){
        console.log(v);
        sAnt = v.vienen == null?0:v.vienen;
        sTotal = v.van == null?0:v.van;
        ingreso = v.ingreso == null?0:v.ingreso;
        egreso = v.egreso == null?0:v.egreso;
        prestamo = v.prestamo == null?0:v.prestamo;
        entregado = v.entregado == null?0:v.entregado;
        recibido = v.recibido == null?0:v.recibido;
        cuota = v.cuota == null?0:v.cuota;
        movimiento = (parseFloat(egreso) + parseFloat(cuota)) - (parseFloat(ingreso) + parseFloat(prestamo)) ;
        saldo = v.saldo == null?0:v.saldo;
        x = parseFloat(entregado) - parseFloat(recibido);
        total = parseFloat(v.saldo) + movimiento + x + sAnt;

        fil = psFila(fila,v.fecha);
        //rS.cell(fil,1).data(sAnt.toFixed(2)).draw();
        rS.cell(fil,2).data(sAnt.toFixed(2)).draw();

        rS.cell(fil,3).data(saldo.toFixed(2)).draw();
        rS.cell(fil,4).data(movimiento.toFixed(2)).draw();
        rS.cell(fil,5).data(entregado.toFixed(2)).draw();
        rS.cell(fil,6).data(recibido.toFixed(2)).draw();

        rS.cell(fil,7).data(sTotal.toFixed(2)).draw();
        //sAnt = total;

    } );


    var table = $('#reporteSaldosGeneral').DataTable();
        $('#reporteSaldosGeneral tbody').on( 'click', 'tr', function () {
            var data = table.row( this ).data();
            var lot = $('#reporteSaldosSistemas').DataTable();
            var par = $('#reporteSaldosSistemasParley').DataTable();
            lot.clear().draw();
            par.clear().draw();
            VentanaEmergente(data[0], data[1]);
        } );

}

function CBTotales(){
    $("#divTabla").html('\
    <table class="table table-striped table-hover table-bordered"  id="reporteSaldosGeneral" width="100%">\
        <thead>\
        <tr>\
            <th>Fecha</th>\
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
    var rS = $('#reporteSaldosGeneral').DataTable(
        {
            "order": [[ 0, "asc" ]],
            "createdRow": function ( row, data, index ) {
                        if ( data[2] * 1 > 15 ) {
                            $('td', row).eq(2).addClass('highlight');
                        }
                    }
        }
    );
    rS.column( 0 ).visible( false );

    rS.clear().draw();
    return rS;
}


/**
 *
 */
function VentanaEmergente(id, fechaAux){
    $('#lbldia').text(id);
    $('#lbldiaA').text(fechaAux);
    $('#ventanaReporteTitulo').html('Saldos por programas para el día (' + fechaAux + ')');
    CargarDatosSistemas(id, $("#agencia").html());
    CargarDatosSistemasParley(id, $("#agencia").html());
    $('#ventanaReporte').modal({ keyboard: false });   // initialized with no keyboard
    $('#ventanaReporte').modal('show');                // initializes and invokes show
}

function CargarDatosSistemas(fecha, agencia){
    var data = JSON.stringify({
        agencia: agencia,
        desde:fecha,
        hasta:fecha
    });
    url = "api/balance/cobrosypagossistemas";
    $.post(url,data)
    .done(function(data){
        var rs = $('#reporteSaldosSistemas').DataTable();
        rs.clear().draw(false);
        rs.row.add([0,0,0,0,0,0,0,0,0]).draw(false);
        suma = 0;
        $.each(data, function(c, v){
            if (v.archivo == null){
                suma += v.saldo;
                var id = SeleccionarSistemas(v.sistema);
                rs.cell(0,id).data(v.saldo.toFixed(2)).draw();
            }
        });
        $('#saldoloteria').html('Total Saldo Por Loteria (' + suma.toFixed(2) + ')');
    });
}

function CargarDatosSistemasParley(fecha, agencia){
    var data = JSON.stringify({
        agencia: agencia,
        desde:fecha,
        hasta:fecha
    });

    url = "api/balance/cobrosypagossistemas";
    $.post(url,data)
    .done(function(data){
        var rs = $('#reporteSaldosSistemasParley').DataTable();
        rs.clear().draw(false);

        rs.row.add([0,0,0]).draw(false);
        suma = 0;
        $.each(data, function(c, v){
            if (v.archivo != null){
                suma += v.saldo;
                var id = SeleccionarSistemas(v.sistema);
                rs.cell(0,id).data(v.saldo.toFixed(2)).draw();

            }
        });
        $('#saldoparley').html('Total Saldo Por Parley (' + suma.toFixed(2) + ')');
    });
}


function SeleccionarSistemas(id){
    switch (id) {
        case 1:
            return 0;
            break;
        case 2:
            return 1;
            break;
        case 3:
            return 2;
            break;
        case 4:
            return 3;
            break;
        case 5:
            return 4;
            break;
        case 6:
            return 0;
            break;
        case 7:
            return 1;
            break;
        case 8:
            return 2;
            break;
        case 9:
            return 5;
            break;
        case 10:
            return 6;
            break;
        case 11:
            return 7;
            break;
        case 12:
            return 8;
            break;
        default:
            break;
    }
}

function VerDetallesTaquillas(id){

    fecha =  $('#lbldia').text();
    fechaA =  $('#lbldiaA').text();
    $('#ventanaEmergenteTitulo').html('Detalles del día (' + fechaA + ')');
    var data = JSON.stringify({
        agencia: $("#agencia").html(),
        desde:fecha,
        hasta:fecha
    });

    url = "api/balance/cobrosypagosdetallados";
    $.post(url,data)
    .done(function(data){
        tabla = '<table class="table table-striped table-bordered table-hover" data-page-length="5" id="reporteSaldosDetallados" width="100%"><thead><tr>\
        <th>Taquila</th><th>Venta</th><th>Premio</th><th>Comision</th><th>Saldo</th><th>Programa</th>\
        <tbody>';
        tabla += '</tbody></tr></thead></table>';
        $('#ventanaEmergenteContenido').html(tabla);
        var rs = $('#reporteSaldosDetallados').DataTable({
            "ordering":  false,
            "info":      false,
            "searching": false,
            "paging":    true
        });
        rs.clear().draw();
        $.each(data, function(c, v){
            venta = v.venta == null?0:v.venta;
            premio = v.premio == null?0:v.premio;
            comision = v.comision == null?0:v.comision;
            saldo = v.saldo == null?0:v.saldo;
            //tabla += '<tr><td>' + v.taquilla + '</td><td>' + venta + '</td><td>' + premio + '</td>\
            //<td>' + comision + '</td><td>' + saldo + '</td><td>' + v.observacion + '</td></tr>'
            rs.row.add([
                v.taquilla,
                venta,
                premio,
                comision,
                saldo,
                v.observacion
            ]).draw(false);
        });



    });

    $('#ventanaEmergente').modal({ keyboard: false });   // initialized with no keyboard
    $('#ventanaEmergente').modal('show');                // initializes and invokes show
}


/**
 * @param Date | UNIX
 * @param Date | UNIX
 * @param DataTable
 */
function RecorreFechas(desde, hasta, rS){
    fauxd = desde.split("/");
    fauxh = hasta.split("/");
    danio = parseInt(fauxd[0]);
    dmes =  parseInt(fauxd[1]);
    ddia =  parseInt(fauxd[2]);

    hanio = parseInt(fauxh[0]);
    hmes =  parseInt(fauxh[1]);
    hdia =  parseInt(fauxh[2]);
    var fila = {};
    var count = 0;
    for (h = danio; h <= hanio; h++){
        for (i = dmes; i <= hmes; i++){
            mdmes = new Date(h, i, 0).getDate();
            for (j = ddia; j <= mdmes; j++){
                dia = j;
                if((String(j)).length==1)dia='0'+j;
                mes = i;
                if((String(i)).length==1)mes='0'+i;
                fecha = h + "-" + mes + "-" + dia;
                fechaX = dia + "/" + mes + "/" + h;
                dia = new Date(h + "/" + mes + "/" + dia);
                pos = dia.getDay();
                inicio = '';
                fin = '';
                if( pos == 0){
                    inicio = '<label style="color:green">';
                    fin = '</label>'
                }
                fila[count] = fecha;
                 rS.row.add( [
                    fecha,
                    inicio + _SAV[pos] + ' ' + fechaX + fin,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                ] ).draw( false );
                count++;
                //console.log(danio + "-" + i + "-" + j);
                if(hanio == h && hmes == i && hdia == j )break;
            }
            ddia = 1;
        }
    }
    return fila;
}

function btnAccion(id){
    s = '<div class="btn-group">\
        <button type="button" class="btn btn-success">\
        Action</button>\
        <button type="button" class="btn btn-success dropdown-toggle" \
        data-toggle="dropdown" aria-expanded="false">\
        <span class="caret"></span>\
        <span class="sr-only">Toggle Dropdown</span>\
        </button>\
        <ul class="dropdown-menu" role="menu">\
        <li><a href="#">Action</a></li>\
        <li><a href="#">Another action</a></li>\
        <li><a href="#">Something else here</a></li>\
        <li class="divider"></li>\
        <li><a href="#">Separated link</a></li>\
        </ul>\
    </div>';
    return s
}

/**
 * Listar Cuentas para Movimientos
 */
function LCuentaM(id){

}
