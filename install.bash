PREFIX="$1"
DESTDIR="$2"

engine_name="bamboo"
engine_gui_name="ibus-setup-Bamboo.desktop"
ibus_e_name="ibus-engine-${engine_name}"
pkg_name="ibus-${engine_name}"
version="0.8.4"

engine_dir=${PREFIX}/share/${pkg_name}
ibus_dir=${PREFIX}/share/ibus

# Main script
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
	mkdir -p ${DESTDIR}${engine_dir}
	mkdir -p ${DESTDIR}${PREFIX}/lib/ibus-${engine_name}
	mkdir -p ${DESTDIR}${ibus_dir}/component/
	mkdir -p ${DESTDIR}${PREFIX}/share/applications/

	cp -R -f icons data ${DESTDIR}${engine_dir}
	cp -f ${ibus_e_name} ${DESTDIR}${PREFIX}/lib/ibus-${engine_name}/
	cp -f data/${engine_name}.xml ${DESTDIR}${ibus_dir}/component/
	cp -f data/${engine_gui_name} ${DESTDIR}${PREFIX}/share/applications/
elif [[ "$OSTYPE" == "freebsd"* ]]; then
	mkdir -p ${DESTDIR}${engine_dir}
	mkdir -p ${DESTDIR}${PREFIX}/lib/ibus-${engine_name}
	mkdir -p ${DESTDIR}${ibus_dir}/component/
	mkdir -p ${DESTDIR}${PREFIX}/share/applications/

	cp -R -f icons data ${DESTDIR}${engine_dir}
	cp -f ${ibus_e_name} ${DESTDIR}${PREFIX}/lib/ibus-${engine_name}/
	cp -f data/${engine_name}-freebsd.xml ${DESTDIR}${ibus_dir}/component/
	cp -f data/${engine_gui_name} ${DESTDIR}${PREFIX}/share/applications/
else
	echo "Operating system is not supported currently!"
fi
